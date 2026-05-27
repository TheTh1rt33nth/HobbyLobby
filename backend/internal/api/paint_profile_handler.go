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

type PaintProfileHandler struct {
	paintProfileStore store.PaintProfileStore
	logger            *log.Logger
}

func NewPaintProfileHandler(paintProfileStore store.PaintProfileStore, logger *log.Logger) *PaintProfileHandler {
	return &PaintProfileHandler{
		paintProfileStore: paintProfileStore,
		logger:            logger,
	}
}

func (pph *PaintProfileHandler) GetPaintProfilesByUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r)

	profiles, err := pph.paintProfileStore.GetPaintProfilesByUserId(r.Context(), user.Id)
	if err != nil {
		pph.logger.Printf("GetPaintProfilesByUser: failed to list profiles: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"paintProfiles": profiles})
}

func (pph *PaintProfileHandler) GetPaintProfileById(w http.ResponseWriter, r *http.Request) {
	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		pph.logger.Printf("GetPaintProfileById: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	profile, err := pph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		pph.logger.Printf("GetPaintProfileById: failed to get profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if profile == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasPaintProfileAccess(profile, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this paint profile"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"paintProfile": profile})
}

func (pph *PaintProfileHandler) CreatePaintProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r)

	var profile store.PaintProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		pph.logger.Printf("CreatePaintProfile: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	profile.UserId = user.Id

	createdProfile, err := pph.paintProfileStore.CreatePaintProfile(r.Context(), &profile)
	if err != nil {
		pph.logger.Printf("CreatePaintProfile: failed to create profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusCreated, handler_utils.Envelope{"paintProfile": createdProfile})
}

func (pph *PaintProfileHandler) UpdatePaintProfile(w http.ResponseWriter, r *http.Request) {
	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		pph.logger.Printf("UpdatePaintProfile: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	existingProfile, err := pph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		pph.logger.Printf("UpdatePaintProfile: failed to get profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingProfile == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasPaintProfileAccess(existingProfile, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this paint profile"})
		return
	}

	var profile store.PaintProfile
	if err = json.NewDecoder(r.Body).Decode(&profile); err != nil {
		pph.logger.Printf("UpdatePaintProfile: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	updatedProfile, err := pph.paintProfileStore.UpdatePaintProfile(r.Context(), profileId, &profile)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
			return
		}
		pph.logger.Printf("UpdatePaintProfile: failed to update profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"paintProfile": updatedProfile})
}

func (pph *PaintProfileHandler) DeletePaintProfile(w http.ResponseWriter, r *http.Request) {
	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		pph.logger.Printf("DeletePaintProfile: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	existingProfile, err := pph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		pph.logger.Printf("DeletePaintProfile: failed to get profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if existingProfile == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasPaintProfileAccess(existingProfile, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this paint profile"})
		return
	}

	if err = pph.paintProfileStore.DeletePaintProfile(r.Context(), profileId); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
			return
		}
		pph.logger.Printf("DeletePaintProfile: failed to delete profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (pph *PaintProfileHandler) AddPaintStep(w http.ResponseWriter, r *http.Request) {
	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		pph.logger.Printf("AddPaintStep: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	profile, err := pph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		pph.logger.Printf("AddPaintStep: failed to get profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if profile == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasPaintProfileAccess(profile, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this paint profile"})
		return
	}

	var step store.PaintStep
	if err = json.NewDecoder(r.Body).Decode(&step); err != nil {
		pph.logger.Printf("AddPaintStep: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	step.PaintProfileId = profileId

	createdStep, err := pph.paintProfileStore.CreatePaintStep(r.Context(), &step)
	if err != nil {
		if errors.Is(err, store.ErrConflict) {
			handler_utils.WriteJSON(w, http.StatusConflict, handler_utils.Envelope{"error": "a step with that order already exists in this profile"})
			return
		}
		pph.logger.Printf("AddPaintStep: failed to create step: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusCreated, handler_utils.Envelope{"paintStep": createdStep})
}

func (pph *PaintProfileHandler) UpdatePaintStep(w http.ResponseWriter, r *http.Request) {
	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		pph.logger.Printf("UpdatePaintStep: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	stepId, err := handler_utils.GetIntParameterFromRequest(r, "stepId")
	if err != nil {
		pph.logger.Printf("UpdatePaintStep: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid step id"})
		return
	}

	profile, err := pph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		pph.logger.Printf("UpdatePaintStep: failed to get profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if profile == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasPaintProfileAccess(profile, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this paint profile"})
		return
	}

	var step store.PaintStep
	if err = json.NewDecoder(r.Body).Decode(&step); err != nil {
		pph.logger.Printf("UpdatePaintStep: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	updatedStep, err := pph.paintProfileStore.UpdatePaintStep(r.Context(), stepId, &step)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint step not found"})
			return
		}
		pph.logger.Printf("UpdatePaintStep: failed to update step: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"paintStep": updatedStep})
}

func (pph *PaintProfileHandler) DeletePaintStep(w http.ResponseWriter, r *http.Request) {
	profileId, err := handler_utils.GetIntParameterFromRequest(r, "profileId")
	if err != nil {
		pph.logger.Printf("DeletePaintStep: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid profile id"})
		return
	}

	stepId, err := handler_utils.GetIntParameterFromRequest(r, "stepId")
	if err != nil {
		pph.logger.Printf("DeletePaintStep: failed to read parameter from URL: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid step id"})
		return
	}

	profile, err := pph.paintProfileStore.GetPaintProfileById(r.Context(), profileId)
	if err != nil {
		pph.logger.Printf("DeletePaintStep: failed to get profile: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}
	if profile == nil {
		handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint profile not found"})
		return
	}

	user := middleware.GetUserContext(r)
	if !auth_utils.UserHasPaintProfileAccess(profile, user) {
		handler_utils.WriteJSON(w, http.StatusForbidden, handler_utils.Envelope{"error": "you do not have access to this paint profile"})
		return
	}

	if err = pph.paintProfileStore.DeletePaintStep(r.Context(), stepId); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			handler_utils.WriteJSON(w, http.StatusNotFound, handler_utils.Envelope{"error": "paint step not found"})
			return
		}
		pph.logger.Printf("DeletePaintStep: failed to delete step: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
