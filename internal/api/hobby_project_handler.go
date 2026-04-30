package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HobbyProjectHandler struct {
}

func NewHobbyProjectHandler() *HobbyProjectHandler {
	return &HobbyProjectHandler{}
}

func (hp *HobbyProjectHandler) GetHobbyProjectById(w http.ResponseWriter, r *http.Request) {
	paramsHoobbyProjectId := chi.URLParam(r, "projectId")
	if paramsHoobbyProjectId == "" {
		http.NotFound(w, r)
		return
	}

	hobbyProjectId, err := strconv.Atoi(paramsHoobbyProjectId)
	if err != nil {
		http.NotFound(w, r) // TODO: handle err properly
		return
	}

	fmt.Fprintf(w, "ProjectId: %d", hobbyProjectId)
}
