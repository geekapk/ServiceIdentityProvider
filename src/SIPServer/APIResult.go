package SIPServer

import (
	"net/http"
	"encoding/json"
)

type APIResult struct {
	Ok bool `json:"ok"`
	Data interface{} `json:"data"`
}

func BuildAPIResult(ok bool, data interface{}) *APIResult {
	return &APIResult {
		Ok: ok,
		Data: data,
	}
}

func (r *APIResult) Write(w http.ResponseWriter) {
	data, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
