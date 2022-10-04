package itface

import (
	"github.com/khuchuz/go-clean-architecture-sql/auth/entities"
	"github.com/khuchuz/go-clean-architecture-sql/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(inp entities.SignUpInput) error
	SignIn(inp entities.SignInput) (string, error)
	ChangePassword(inp entities.ChangePasswordInput) error
	ParseToken(accessToken string) (*models.User, error)
	DeleteAccount(inp entities.DeleteInput) error
}
