package services

import (
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(inp models.SignUpInput) error
	SignIn(inp models.SignInput) (string, error)
	ChangePassword(inp models.ChangePasswordInput) error
	ParseToken(accessToken string) (*models.User, error)
	DeleteAccount(inp models.DeleteInput) error
}
