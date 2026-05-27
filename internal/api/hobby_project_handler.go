package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/middleware"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/auth_utils"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/handler_utils"
)

type HobbyProjectHandler struct {
	projectStore      store.HobbyProjectStore
	paintProfileStore store.PaintProfileStore
	unitStore         store.UnitStore
	logger            *log.Logger
}

func NewHobbyProjectHandler(projectStore store.HobbyProjectStore, paintProfileStore store.PaintProfileStore, unitStore store.UnitStore, logger *log.Logger) *HobbyProjectHandler {
	return &HobbyProjectHandler{
		projectStore:      projectStore,
		paintProfileStore: paintProfileStore,
		unitStore:         unitStore,
		logger:            logger,
	}
}

func (hph *HobbyProjectHandler) GetHobbyProjectsByUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r)

	projects, err := hph.projectStore.GetHobbyProjectsByUserId(r.Context(), user.Id)
	if err != nil {
		hph.logger.Printf("GetHobbyProjectsByUser: failed to list projects: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	projectIds := make([]int, len(projects))
	for i, p := range projects {
		projectIds[i] = p.Id
	}
	progressMap, err := hph.unitStore.GetProgressByProjectIds(r.Context(), projectIds)
	if err != nil {
		hph.logger.Printf("GetHobbyProjectsByUser: failed to get progress map: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	for _, p := range projects {
		if prog, ok := progressMap[p.Id]; ok {
			p.Progress = prog
		} else {
			p.Progress = &store.ProjectProgress{ByStatus: make(map[string]int)}
		}
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"hobbyProjects": projects})
}

func (hph *HobbyProjectHandler) GetHobbyProjectById(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("GetHobbyProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	project, err := hph.projectStore.GetHobbyProjectById(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("GetHobbyProjectById: failed to get HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if project == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasAccess(project, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this project"})
		return
	}

	progress, err := hph.unitStore.GetProgressByProjectId(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("GetHobbyProjectById: failed to get progress: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	project.Progress = progress

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"hobbyProject": project})
}

func (hph *HobbyProjectHandler) CreateHobbyProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r)

	var project store.HobbyProject
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		hph.logger.Printf("CreateHobbyProject: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	project.UserId = user.Id

	createdProject, err := hph.projectStore.CreateHobbyProject(r.Context(), &project)
	if err != nil {
		if errors.Is(err, store.ErrInvalidInput) {
			handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid field value"})
			return
		}
		hph.logger.Printf("CreateHobbyProject: failed to create HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusCreated, handler_utils.Envelope{"hobbyProject": createdProject})
}

func (hph *HobbyProjectHandler) UpdateHobbyProject(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	existingProject, err := hph.projectStore.GetHobbyProjectById(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to get existing HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingProject == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasAccess(existingProject, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this project"})
		return
	}

	var project store.HobbyProject
	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		hph.logger.Printf("UpdateHobbyProject: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	updatedProject, err := hph.projectStore.UpdateHobbyProject(r.Context(), hobbyProjectId, &project)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
			return
		}
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
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	existingProject, err := hph.projectStore.GetHobbyProjectById(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("DeleteHobbyProject: failed to get existing HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingProject == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasAccess(existingProject, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this project"})
		return
	}

	if err = hph.projectStore.DeleteHobbyProject(r.Context(), existingProject.Id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
			return
		}
		hph.logger.Printf("DeleteHobbyProject: failed to delete HobbyProject: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (hph *HobbyProjectHandler) GetPaintProfilesForProject(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("GetPaintProfilesForProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	project, err := hph.projectStore.GetHobbyProjectById(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("GetPaintProfilesForProject: failed to get project: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if project == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasAccess(project, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this project"})
		return
	}

	profiles, err := hph.paintProfileStore.GetPaintProfilesByProjectId(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("GetPaintProfilesForProject: failed to list paint profiles: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"paintProfiles": profiles})
}

func (hph *HobbyProjectHandler) AssignPaintProfile(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("AssignPaintProfile: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		hph.logger.Printf("AssignPaintProfile: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	user := middleware.GetUserContext(r)

	project, err := hph.projectStore.GetHobbyProjectById(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("AssignPaintProfile: failed to get project: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if project == nil || !auth_utils.UserHasAccess(project, user) {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
		return
	}

	profile, err := hph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		hph.logger.Printf("AssignPaintProfile: failed to get paint profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if profile == nil || !auth_utils.UserHasPaintProfileAccess(profile, user) {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	if err = hph.paintProfileStore.AssignToProject(r.Context(), hobbyProjectId, profileId); err != nil {
		hph.logger.Printf("AssignPaintProfile: failed to assign: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (hph *HobbyProjectHandler) UnassignPaintProfile(w http.ResponseWriter, r *http.Request) {
	hobbyProjectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		hph.logger.Printf("UnassignPaintProfile: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		hph.logger.Printf("UnassignPaintProfile: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	user := middleware.GetUserContext(r)

	project, err := hph.projectStore.GetHobbyProjectById(r.Context(), hobbyProjectId)
	if err != nil {
		hph.logger.Printf("UnassignPaintProfile: failed to get project: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if project == nil || !auth_utils.UserHasAccess(project, user) {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "project not found"})
		return
	}

	if err = hph.paintProfileStore.UnassignFromProject(r.Context(), hobbyProjectId, profileId); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not assigned to this project"})
			return
		}
		hph.logger.Printf("UnassignPaintProfile: failed to unassign: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
