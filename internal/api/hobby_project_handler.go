package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/go-chi/chi/v5"
)

type HobbyProjectHandler struct {
	projectStore store.HobbyProjectStore
}

func NewHobbyProjectHandler(projectStore store.HobbyProjectStore) *HobbyProjectHandler {
	return &HobbyProjectHandler{
		projectStore: projectStore,
	}
}

func (hph *HobbyProjectHandler) GetHobbyProjectById(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := getProjectIdFromRequest(r)
	if err != nil {
		http.NotFound(w, r) // TODO: handle err properly
		return
	}

	project, err := hph.projectStore.GetHobbyProjectById(hobbyProjectId)
	if err != nil {
		http.Error(w, "Failed to get project", http.StatusInternalServerError)
		return
	}

	addDefaultResponseHeaders(w)
	json.NewEncoder(w).Encode(project)
}

func (hph *HobbyProjectHandler) CreateHobbyProject(w http.ResponseWriter, r *http.Request) {
	var project store.HobbyProject
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Failed to create project - invalid payload", http.StatusBadRequest)
		return
	}

	createdProject, err := hph.projectStore.CreateHobbyProject(&project)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	addDefaultResponseHeaders(w)
	json.NewEncoder(w).Encode(createdProject)

}

func (hph *HobbyProjectHandler) UpdateHobbyProject(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := getProjectIdFromRequest(r)
	if err != nil {
		http.NotFound(w, r) // TODO: handle err properly
		return
	}

	existingProject, err := hph.projectStore.GetHobbyProjectById(hobbyProjectId)
	if err != nil {
		http.Error(w, "Failed to update project - failed to get existing project", http.StatusInternalServerError)
		return
	}
	if existingProject == nil {
		http.NotFound(w, r)
		return
	}

	var project store.HobbyProject
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Failed to update project - invalid payload", http.StatusBadRequest)
		return
	}

	updatedProject, err := hph.projectStore.UpdateHobbyProject(hobbyProjectId, &project)
	if err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	addDefaultResponseHeaders(w)
	json.NewEncoder(w).Encode(updatedProject)
}

func (hph *HobbyProjectHandler) DeleteHobbyProject(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := getProjectIdFromRequest(r)
	if err != nil {
		http.NotFound(w, r) // TODO: handle err properly
		return
	}

	err = hph.projectStore.DeleteHobbyProject(hobbyProjectId)
	if err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getProjectIdFromRequest(r *http.Request) (int, error) {
	paramsHobbyProjectId := chi.URLParam(r, "projectId")
	if paramsHobbyProjectId == "" {
		return -1, fmt.Errorf("projectId param is required")
	}

	hobbyProjectId, err := strconv.Atoi(paramsHobbyProjectId)
	if err != nil {
		return -1, err
	}

	return hobbyProjectId, nil
}

func addDefaultResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
