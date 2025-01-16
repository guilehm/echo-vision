package web

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, data any) {
	_ = json.NewEncoder(w).Encode(data)
}
