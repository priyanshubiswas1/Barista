package menu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	menupkg "github.com/priyanshubiswas1/barista/menuservice/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ScanWorker struct {
	repo     IRepository
	uploader *menupkg.GarageUploader

	python string
	script string
}

func NewScanWorker(repo IRepository, uploader *menupkg.GarageUploader, python, script string) *ScanWorker {
	if python == "" {
		python = "python"
	}
	return &ScanWorker{repo: repo, uploader: uploader, python: python, script: script}
}

func (w *ScanWorker) HandleMessage(ctx context.Context, body []byte) error {
	if w == nil || w.repo == nil || w.uploader == nil {
		return errors.New("scan worker not configured")
	}

	var job ScanJob
	if err := json.Unmarshal(body, &job); err != nil {
		return err
	}
	if job.ItemID == "" || job.ShopID == "" || job.ImageURL == "" {
		return errors.New("invalid scan job payload")
	}
	log.Println("processing 3d scan job", "shopId", job.ShopID, "itemId", job.ItemID)

	imgPath, cleanup, err := downloadToTemp(ctx, job.ImageURL)
	if err != nil {
		return err
	}
	defer cleanup()

	out, err := w.runHunyuan(ctx, imgPath)
	if err != nil {
		return err
	}
	out, err = extractFirstJSONObject(out)
	if err != nil {
		return fmt.Errorf("hunyuan output did not contain json: %w", err)
	}

	// The python script prints {"result": <something>}. We store that json blob as a .json artifact
	// AND attempt to upload any referenced file path if present.
	artifactKeyBase := fmt.Sprintf("hunyuan/%s/%s", job.ShopID, job.ItemID)

	// Prefer the best 3D artifact (GLB > GLTF > ZIP > OBJ).
	data, contentType, ext, err := extractBestArtifactBytes(out)
	if err == nil && len(data) > 0 {
		filename := "model" + ext
		key := fmt.Sprintf("%s/%s", artifactKeyBase, filename)
		url, upErr := w.uploader.UploadBytes(ctx, key, contentType, data)
		if upErr != nil {
			return upErr
		}
		_, updErr := w.repo.SetItemScan3DURL(ctx, job.ItemID, job.ShopID, url, time.Now())
		if updErr == nil {
			log.Println("stored 3d scan url", "shopId", job.ShopID, "itemId", job.ItemID, "url", url)
		}
		return updErr
	}
	if err != nil {
		snippet := string(out)
		if len(snippet) > 2000 {
			snippet = snippet[:2000] + "...(truncated)"
		}
		log.Println("no 3d artifact found in hunyuan output; uploading debug json", "shopId", job.ShopID, "itemId", job.ItemID, "err", err, "output", snippet)
	}

	// Fallback: upload the raw json output for debugging and don't update scan url.
	key := fmt.Sprintf("%s/%s", artifactKeyBase, "hunyuan_result.json")
	url, upErr := w.uploader.UploadBytes(ctx, key, "application/json", out)
	if upErr == nil {
		log.Println("uploaded hunyuan debug output", "shopId", job.ShopID, "itemId", job.ItemID, "url", url)
	}
	return upErr
}

func (w *ScanWorker) runHunyuan(ctx context.Context, imagePath string) ([]byte, error) {
	if strings.TrimSpace(w.script) == "" {
		return nil, errors.New("HUNYUAN_SCRIPT not configured")
	}
	cmd := exec.CommandContext(ctx, w.python, w.script, "--image", imagePath)
	cmd.Env = os.Environ()
	return cmd.CombinedOutput()
}

func downloadToTemp(ctx context.Context, url string) (string, func(), error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", func() {}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", func() {}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", func() {}, fmt.Errorf("failed to fetch image: status %s", resp.Status)
	}

	ext := ".img"
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "png") {
		ext = ".png"
	} else if strings.Contains(ct, "jpeg") || strings.Contains(ct, "jpg") {
		ext = ".jpg"
	}

	tmpDir := os.TempDir()
	name := uuid.New().String() + ext
	path := filepath.Join(tmpDir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", func() {}, err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(path)
		return "", func() {}, err
	}

	return path, func() { _ = os.Remove(path) }, nil
}

func extractBestArtifactBytes(hunyuanJSON []byte) ([]byte, string, string, error) {
	// The output shape varies. We'll unmarshal and search recursively for a string
	// that is either a local path or a downloadable URL.
	var v any
	if err := json.Unmarshal(hunyuanJSON, &v); err != nil {
		return nil, "", "", err
	}
	refs := findAllArtifactRefs(v, nil)
	ref := pickBestArtifactRef(refs)
	if ref == "" {
		return nil, "", "", errors.New("no artifact reference found in hunyuan output")
	}

	// If it's a URL, download it. Otherwise treat as a local path.
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		req, err := http.NewRequest(http.MethodGet, ref, nil)
		if err != nil {
			return nil, "", "", err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, "", "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, "", "", fmt.Errorf("failed to fetch artifact: status %s", resp.Status)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", "", err
		}
		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			ct = guessContentTypeFromName(ref)
		}
		return b, ct, guessExtFromName(ref), nil
	}

	b, err := os.ReadFile(ref)
	if err != nil {
		return nil, "", "", err
	}
	return b, guessContentTypeFromName(ref), guessExtFromName(ref), nil
}

func extractFirstJSONObject(out []byte) ([]byte, error) {
	// Some clients (gradio_client) may print informational lines before JSON.
	// We locate the first '{' and last '}' and parse that slice.
	start := bytes.IndexByte(out, '{')
	end := bytes.LastIndexByte(out, '}')
	if start < 0 || end < 0 || end <= start {
		return nil, errors.New("no json object found in output")
	}
	candidate := bytes.TrimSpace(out[start : end+1])

	// Validate it is JSON.
	var tmp any
	if err := json.Unmarshal(candidate, &tmp); err != nil {
		return nil, err
	}
	return candidate, nil
}

func findAllArtifactRefs(v any, out []string) []string {
	switch t := v.(type) {
	case map[string]any:
		// Common shapes include objects with "url" or "path" keys.
		if u, ok := t["url"].(string); ok {
			s := strings.TrimSpace(u)
			if s != "" {
				out = append(out, s)
			}
		}
		if p, ok := t["path"].(string); ok {
			s := strings.TrimSpace(p)
			if s != "" {
				out = append(out, s)
			}
		}
		if v, ok := t["value"].(string); ok {
			s := strings.TrimSpace(v)
			if s != "" {
				out = append(out, s)
			}
		}
		for _, vv := range t {
			out = findAllArtifactRefs(vv, out)
		}
	case []any:
		for _, vv := range t {
			out = findAllArtifactRefs(vv, out)
		}
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return out
		}
		if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
			out = append(out, s)
			return out
		}
		if filepath.IsAbs(s) {
			out = append(out, s)
			return out
		}
	}
	return out
}

func guessContentTypeFromName(name string) string {
	l := strings.ToLower(name)
	switch {
	case strings.Contains(l, ".glb"):
		return "model/gltf-binary"
	case strings.Contains(l, ".gltf"):
		return "model/gltf+json"
	case strings.Contains(l, ".zip"):
		return "application/zip"
	case strings.Contains(l, ".obj"):
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

func guessExtFromName(name string) string {
	l := strings.ToLower(name)
	switch {
	case strings.Contains(l, ".glb"):
		return ".glb"
	case strings.Contains(l, ".gltf"):
		return ".gltf"
	case strings.Contains(l, ".zip"):
		return ".zip"
	case strings.Contains(l, ".obj"):
		return ".obj"
	default:
		return ".bin"
	}
}

func pickBestArtifactRef(refs []string) string {
	best := ""
	bestScore := -1
	for _, r := range refs {
		l := strings.ToLower(r)
		score := 0
		switch {
		case strings.Contains(l, ".textured_mesh.glb") || strings.HasSuffix(l, "textured_mesh.glb"):
			score = 400
		case strings.Contains(l, ".glb"):
			score = 300
		case strings.Contains(l, ".gltf"):
			score = 200
		case strings.Contains(l, ".zip"):
			score = 150
		case strings.Contains(l, ".obj"):
			score = 100
		default:
			score = 0
		}
		if score > bestScore {
			bestScore = score
			best = r
		}
	}
	return best
}

func StartScanConsumer(
	repo IRepository,
	uploader *menupkg.GarageUploader,
	conn *amqp.Connection,
	exchange, queueName, bindingKey, consumerName string,
	python, script string,
) error {
	if conn == nil {
		return errors.New("rabbitmq connection is nil")
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, bindingKey, exchange, false, nil); err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, consumerName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	worker := NewScanWorker(repo, uploader, python, script)
	go func() {
		for msg := range msgs {
			go func(m amqp.Delivery) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
				defer cancel()

				if err := worker.HandleMessage(ctx, m.Body); err != nil {
					log.Println("scan job failed:", err)
					_ = m.Nack(false, false) // dead-letter behavior not configured; don't requeue by default
					return
				}
				log.Println("scan job completed")
				_ = m.Ack(false)
			}(msg)
		}
	}()

	return nil
}

