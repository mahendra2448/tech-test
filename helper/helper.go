package helper

import (
	"encoding/json"
	"net/http"
)

type RespJSON struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Errors  *error `json:"errors"`
}

func ReturnOK(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RespJSON{
		Status:  http.StatusText(http.StatusOK),
		Message: msg,
		Errors:  nil,
	})
}

func ReturnBadRequest(w http.ResponseWriter, err error, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(RespJSON{
		Status:  http.StatusText(http.StatusBadRequest),
		Message: msg,
		Errors:  &err,
	})
}

func ReturnErr(w http.ResponseWriter, err error, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(RespJSON{
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: msg,
		Errors:  &err,
	})
}
