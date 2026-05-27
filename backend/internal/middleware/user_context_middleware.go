package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/tokens"
	"github.com/TheTh1rt33nth/HobbyLobby/internal/utils/handler_utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUserContext(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)

	return r.WithContext(ctx)
}

func GetUserContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		return nil
	}

	return user
}

func (um *UserMiddleware) PopulateUserContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUserContext(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			handler_utils.WriteJSON(w, http.StatusUnauthorized, handler_utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserByToken(r.Context(), tokens.ScopeAuthentication, token)
		if err != nil {
			handler_utils.WriteJSON(w, http.StatusInternalServerError, handler_utils.Envelope{"error": "internal server error"})
			return
		}
		if user == nil {
			handler_utils.WriteJSON(w, http.StatusUnauthorized, handler_utils.Envelope{"error": "invalid token"})
			return
		}

		r = SetUserContext(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireAuthenticatedUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserContext(r)
		if user == nil || user.IsAnonymous() {
			handler_utils.WriteJSON(w, http.StatusUnauthorized, handler_utils.Envelope{"error": "authentication required for this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
