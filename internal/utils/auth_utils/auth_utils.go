package auth_utils

import "github.com/TheTh1rt33nth/HobbyLobby/internal/store"

func UserHasAccess(hp *store.HobbyProject, user *store.User) bool {
	return !user.IsAnonymous() && hp.UserId == user.Id
}

func UserHasPaintProfileAccess(pp *store.PaintProfile, user *store.User) bool {
	return !user.IsAnonymous() && pp.UserId == user.Id
}
