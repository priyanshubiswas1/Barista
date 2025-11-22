package pkg

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type CloudinaryService struct {
	apiKey    string
	apiSecret string
	cloudName string
}

type CloudinaryUploadResponse struct {
	PublicId string `json:"public_id"`
	Url      string `json:"secure_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format  string `json:"format"`
}

func NewCloudinaryService(apiKey, apiSecret, cloudName string) *CloudinaryService {
	return &CloudinaryService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		cloudName: cloudName,
	}
}

func (c *CloudinaryService) UploadImage(ctx context.Context, imageData []byte, filename string) (string, error) {
	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Write image data to the form
	_, err = part.Write(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	// Add timestamp for signed uploads
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	writer.WriteField("timestamp", timestamp)
	writer.WriteField("api_key", c.apiKey)

	// Calculate signature for signed uploads
	if c.apiSecret != "" {
		// Create signature string: timestamp + api_secret
		signatureString := fmt.Sprintf("timestamp=%s%s", timestamp, c.apiSecret)
		hash := sha1.Sum([]byte(signatureString))
		signature := fmt.Sprintf("%x", hash)
		writer.WriteField("signature", signature)
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the HTTP request
	uploadURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", c.cloudName)
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("cloudinary upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var uploadResp CloudinaryUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return uploadResp.Url, nil
}

