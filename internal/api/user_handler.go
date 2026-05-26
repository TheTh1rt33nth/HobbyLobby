package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/dto"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/handler_utils"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/validation_utils"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("RegisterUser: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	err = validation_utils.ValidateRegisterUserRequest(&req)
	if err != nil {
		uh.logger.Printf("RegisterUser: request validation failed: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	err = user.Password.Set(req.Password)
	if err != nil {
		uh.logger.Printf("RegisterUser: failed to hash password: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	user, err = uh.userStore.CreateUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, store.ErrConflict) {
			handler_utils.WriteJSON(w, http.StatusConflict, handler_utils.Envelope{"error": "username or email already in use"})
			return
		}
		uh.logger.Printf("RegisterUser: failed to create user: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusCreated, handler_utils.Envelope{"user": user})

}
