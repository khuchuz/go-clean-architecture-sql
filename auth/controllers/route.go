package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/khuchuz/go-clean-architecture-sql/auth/services"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc services.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
		authEndpoints.POST("/change-pass", h.ChangePassword)
		authEndpoints.POST("/delete-me", h.DeleteAccount)
	}
}
