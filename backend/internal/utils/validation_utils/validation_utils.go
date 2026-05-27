package validation_utils

import (
	"errors"
	"regexp"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/dto"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// TODO: return array of errors instead of just the first one encountered
func ValidateRegisterUserRequest(req *dto.RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
