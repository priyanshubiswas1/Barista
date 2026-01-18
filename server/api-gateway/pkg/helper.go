package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, s int, data any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(data)
}

func StripScheme(url string) string {
	if len(url) > 7 && url[:7] == "http://" {
		return url[7:]
	}
	if len(url) > 8 && url[:8] == "https://" {
		return url[8:]
	}
	return url
}
