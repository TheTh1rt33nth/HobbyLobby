package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/dto"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/handler_utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

// TODO: rate limit
func (th *TokenHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		th.logger.Printf("CreateToken: failed to decode payload: %v", err)
		handler_utils.WriteJSON(w, http.StatusBadRequest, handler_utils.Envelope{"error": "invalid payload"})
		return
	}

	user, err := th.userStore.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		th.logger.Printf("CreateToken: failed to fetch user by username: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	if user == nil {
		th.logger.Printf("CreateToken: user not found")
		handler_utils.WriteJSON(w, http.StatusUnauthorized, handler_utils.Envelope{"error": "invalid credentials"})
		return
	}

	passwordMatches, err := user.Password.Matches(req.Password)
	if err != nil {
		th.logger.Printf("CreateToken: failed to compare password hash: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordMatches {
		th.logger.Printf("CreateToken: password does not match")
		handler_utils.WriteJSON(w, http.StatusUnauthorized, handler_utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := th.tokenStore.DeleteAndCreateToken(r.Context(), user.Id, 24*time.Hour, "authentication")
	if err != nil {
		th.logger.Printf("CreateToken: failed to create token: %v", err)
		handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
		return
	}

	handler_utils.WriteJSON(w, http.StatusOK, handler_utils.Envelope{"auth_token": token})

}
