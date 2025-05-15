package responses

import (
	"encoding/json"
	"net/http"
)

type dictionary map[string]interface{}

func HTTPResponse(w http.ResponseWriter, statusCode int, key string, value interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dictionary{
		key: value,
	})
}

func HTTPError(w http.ResponseWriter, err string, statusCode int) {
	http.Error(w, err, statusCode)
}