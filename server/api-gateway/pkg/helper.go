package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, s int, data any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(data)
}
