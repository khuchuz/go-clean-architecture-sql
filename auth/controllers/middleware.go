package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khuchuz/go-clean-architecture-sql/auth"
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"github.com/khuchuz/go-clean-architecture-sql/auth/services"
)

type AuthMiddleware struct {
	usecase services.UseCase
}

func NewAuthMiddleware(usecase services.UseCase) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUnauthorized.Error()})
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUnauthorized.Error()})
		return
	}

	if headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUnauthorized.Error()})
		return
	}

	user, err := m.usecase.ParseToken(headerParts[1])
	if err != nil {
		status := http.StatusInternalServerError
		if err == auth.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		c.JSON(status, models.SignResponse{Message: auth.ErrUnknown.Error()})
		return
	}

	c.Set(services.CtxUserKey, user)
}
