package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RestError struct {
	Error string `json:"error,"`
}

func ProcessInternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(RestError{Error: message})
}

func ProcessBadRequestError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(RestError{Error: message})
}

func ProcessAlreadyExistsError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(RestError{Error: message})
}

func ProcessUnauthorizedError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(RestError{Error: message})
}

func ProcessError(w http.ResponseWriter, err error) {
	var internalError *InternalError
	if err != nil {
		if ok := errors.As(err, &internalError); ok {
			if internalError.Code == http.StatusInternalServerError {
			} else {
			}
			w.WriteHeader(internalError.Code)
			json.NewEncoder(w).Encode(RestError{Error: internalError.Message})
			return
		}
		return
	}
}
