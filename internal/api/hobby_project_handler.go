package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/handler_utils"
)

type HobbyProjectHandler struct {
	projectStore store.HobbyProjectStore
	logger       *log.Logger
}

func NewHobbyProjectHandler(projectStore store.HobbyProjectStore, logger *log.Logger) *HobbyProjectHandler {
	return &HobbyProjectHandler{
		projectStore: projectStore,
		logger:       logger,
	}
}

func (hph *HobbyProjectHandler) GetHobbyProjectById(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("GetHobbyProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid workout Id"})
		return
	}

	project, err := hph.projectStore.GetHobbyProjectById(hobbyProjectId)
	if err != nil {
		hph.logger.Printf("GetHobbyProjectById: failed to update HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"hobbyProject": project})
}

func (hph *HobbyProjectHandler) CreateHobbyProject(w http.ResponseWriter, r *http.Request) {
	var project store.HobbyProject
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		hph.logger.Printf("CreateHobbyProject: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	createdProject, err := hph.projectStore.CreateHobbyProject(&project)
	if err != nil {
		hph.logger.Printf("CreateHobbyProject: failed to update HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"hobbyProject": createdProject})
}

func (hph *HobbyProjectHandler) UpdateHobbyProject(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid workout Id"})
		return
	}

	existingProject, err := hph.projectStore.GetHobbyProjectById(hobbyProjectId)
	if err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to get existing HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingProject == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "HobbyProject not found"})
		return
	}

	var project store.HobbyProject
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	updatedProject, err := hph.projectStore.UpdateHobbyProject(hobbyProjectId, &project)
	if err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to update HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"hobbyProject": updatedProject})
}

func (hph *HobbyProjectHandler) DeleteHobbyProject(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("DeleteHobbyProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid workout Id"})
		return
	}

	err = hph.projectStore.DeleteHobbyProject(hobbyProjectId)
	if err != nil {
		hph.logger.Printf("DeleteHobbyProject: failed to delete HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
