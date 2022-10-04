package delivery

import (
	"github.com/gin-gonic/gin"
	itface "github.com/khuchuz/go-clean-architecture-sql/auth/services/itface"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc itface.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
		authEndpoints.POST("/change-pass", h.ChangePassword)
		authEndpoints.POST("/delete-me", h.DeleteAccount)
	}
}
