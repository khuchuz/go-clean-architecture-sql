package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khuchuz/go-clean-architecture-sql/auth"
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"github.com/khuchuz/go-clean-architecture-sql/auth/services"
	"gorm.io/gorm"
)

type Handler struct {
	useCase services.UseCase
}

func NewHandler(useCase services.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(models.SignUpInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, models.SignResponse{Message: auth.ErrBadRequest.Error()})
		return
	}

	if err := h.useCase.SignUp(*inp); err != nil {
		c.JSON(http.StatusInternalServerError, models.SignResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SignResponse{Message: "Sign Up Berhasil"})
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(models.SignInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, models.SignResponse{Message: auth.ErrBadRequest.Error()})
		return
	}

	token, err := h.useCase.SignIn(*inp)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUnknown.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SignInResponse{Token: token})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	inp := new(models.ChangePasswordInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, models.SignResponse{Message: auth.ErrBadRequest.Error()})
		return
	}

	err := h.useCase.ChangePassword(*inp)
	if err != nil {
		if err == gorm.ErrInvalidTransaction {
			c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrInvalidCreds.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUnknown.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SignResponse{Message: "Password berhasil diubah"})
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	inp := new(models.DeleteInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, models.SignResponse{Message: auth.ErrBadRequest.Error()})
		return
	}

	err := h.useCase.DeleteAccount(*inp)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusUnauthorized, models.SignResponse{Message: auth.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, models.SignResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SignResponse{Message: "Akun berhasil dihapus"})
}
