package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RestError struct {
	Error string `json:"error,"`
}

func ProcessInternalServerError(w http.ResponseWriter, message string) {
	log.Printf("internal server error: %s", message)
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

	if ok := errors.As(err, &internalError); ok {
		if internalError.Code == http.StatusInternalServerError {
			log.Printf("internal server error: %v", err)
		}
		w.WriteHeader(internalError.Code)
		json.NewEncoder(w).Encode(RestError{Error: internalError.Message})
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		log.Printf("error deadline exceeded: %v", err)
		return
	}
	if errors.Is(err, context.Canceled) {
		log.Printf("error context canceled: %v", err)
		return
	}

	log.Printf("unknown error: %v", err)

}
