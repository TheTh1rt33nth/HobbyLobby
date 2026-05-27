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

type UnitHandler struct {
	unitStore    store.UnitStore
	projectStore store.HobbyProjectStore
	logger       *log.Logger
}

func NewUnitHandler(unitStore store.UnitStore, projectStore store.HobbyProjectStore, logger *log.Logger) *UnitHandler {
	return &UnitHandler{
		unitStore:    unitStore,
		projectStore: projectStore,
		logger:       logger,
	}
}

func (uh *UnitHandler) GetUnitsByProject(w http.ResponseWriter, r *http.Request) {
	projectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		uh.logger.Printf("GetUnitsByProject: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	project, err := uh.projectStore.GetHobbyProjectById(r.Context(), projectId)
	if err != nil {
		uh.logger.Printf("GetUnitsByProject: failed to get project: %v", err)
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

	units, err := uh.unitStore.GetUnitsByProjectId(r.Context(), projectId)
	if err != nil {
		uh.logger.Printf("GetUnitsByProject: failed to list units: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"units": units})
}

func (uh *UnitHandler) GetUnitById(w http.ResponseWriter, r *http.Request) {
	projectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		uh.logger.Printf("GetUnitById: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	unitId, err := handler_utils.GetIntParameterFromRequest(r, "unitId")
	if err != nil {
		uh.logger.Printf("GetUnitById: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid unit id"})
		return
	}

	project, err := uh.projectStore.GetHobbyProjectById(r.Context(), projectId)
	if err != nil {
		uh.logger.Printf("GetUnitById: failed to get project: %v", err)
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

	unit, err := uh.unitStore.GetUnitById(r.Context(), unitId)
	if err != nil {
		uh.logger.Printf("GetUnitById: failed to get unit: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if unit == nil || unit.ProjectId != projectId {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "unit not found"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"unit": unit})
}

func (uh *UnitHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	projectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		uh.logger.Printf("CreateUnit: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	project, err := uh.projectStore.GetHobbyProjectById(r.Context(), projectId)
	if err != nil {
		uh.logger.Printf("CreateUnit: failed to get project: %v", err)
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

	var unit store.Unit
	if err = json.NewDecoder(r.Body).Decode(&unit); err != nil {
		uh.logger.Printf("CreateUnit: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	unit.ProjectId = projectId

	createdUnit, err := uh.unitStore.CreateUnit(r.Context(), &unit)
	if err != nil {
		if errors.Is(err, store.ErrInvalidInput) {
			handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid field value"})
			return
		}
		if errors.Is(err, store.ErrForeignKey) {
			handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "paint profile not found"})
			return
		}
		uh.logger.Printf("CreateUnit: failed to create unit: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusCreated, handler_utils.Envelope{"unit": createdUnit})
}

func (uh *UnitHandler) UpdateUnit(w http.ResponseWriter, r *http.Request) {
	projectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		uh.logger.Printf("UpdateUnit: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	unitId, err := handler_utils.GetIntParameterFromRequest(r, "unitId")
	if err != nil {
		uh.logger.Printf("UpdateUnit: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid unit id"})
		return
	}

	project, err := uh.projectStore.GetHobbyProjectById(r.Context(), projectId)
	if err != nil {
		uh.logger.Printf("UpdateUnit: failed to get project: %v", err)
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

	existingUnit, err := uh.unitStore.GetUnitById(r.Context(), unitId)
	if err != nil {
		uh.logger.Printf("UpdateUnit: failed to get unit: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingUnit == nil || existingUnit.ProjectId != projectId {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "unit not found"})
		return
	}

	var unit store.Unit
	if err = json.NewDecoder(r.Body).Decode(&unit); err != nil {
		uh.logger.Printf("UpdateUnit: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	updatedUnit, err := uh.unitStore.UpdateUnit(r.Context(), unitId, &unit)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "unit not found"})
			return
		}
		if errors.Is(err, store.ErrInvalidInput) {
			handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid field value"})
			return
		}
		if errors.Is(err, store.ErrForeignKey) {
			handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "paint profile not found"})
			return
		}
		uh.logger.Printf("UpdateUnit: failed to update unit: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"unit": updatedUnit})
}

func (uh *UnitHandler) DeleteUnit(w http.ResponseWriter, r *http.Request) {
	projectId, err := handler_utils.GetIntParameterFromRequest(r, "projectId")
	if err != nil {
		uh.logger.Printf("DeleteUnit: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid project id"})
		return
	}

	unitId, err := handler_utils.GetIntParameterFromRequest(r, "unitId")
	if err != nil {
		uh.logger.Printf("DeleteUnit: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid unit id"})
		return
	}

	project, err := uh.projectStore.GetHobbyProjectById(r.Context(), projectId)
	if err != nil {
		uh.logger.Printf("DeleteUnit: failed to get project: %v", err)
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

	existingUnit, err := uh.unitStore.GetUnitById(r.Context(), unitId)
	if err != nil {
		uh.logger.Printf("DeleteUnit: failed to get unit: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingUnit == nil || existingUnit.ProjectId != projectId {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "unit not found"})
		return
	}

	if err = uh.unitStore.DeleteUnit(r.Context(), unitId); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "unit not found"})
			return
		}
		uh.logger.Printf("DeleteUnit: failed to delete unit: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
