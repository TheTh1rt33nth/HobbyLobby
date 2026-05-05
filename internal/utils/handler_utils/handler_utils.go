package handler_utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	addDefaultResponseHeaders(w)

	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func GetIntParameterFromRequest(r *http.Request, param string) (int, error) {
	paramValue := chi.URLParam(r, param)
	if paramValue == "" {
		return -1, fmt.Errorf("%s param is required", param)
	}

	id, err := strconv.Atoi(paramValue)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func addDefaultResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
