package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/khuchuz/go-clean-architecture-sql/auth"
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"github.com/khuchuz/go-clean-architecture-sql/auth/services/usecase/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSignUp_Success_200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	signUpBody := &models.SignUpInput{
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "testpass",
	}

	body, err := json.Marshal(signUpBody)
	assert.NoError(t, err)

	uc.On("SignUp", signUpBody.Username, signUpBody.Email, signUpBody.Password).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestSignUp_Failed_400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	body, err := json.Marshal("not json")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
func TestSignUp_Failed_500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	signUpBody := &models.SignUpInput{
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "testpass",
	}

	body, err := json.Marshal(signUpBody)
	assert.NoError(t, err)

	uc.On("SignUp", signUpBody.Username, signUpBody.Email, signUpBody.Password).Return(errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}
func TestSignIn_Sucess_200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	signInBody := &models.SignInput{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(signInBody)
	assert.NoError(t, err)

	uc.On("SignIn", signInBody.Username, signInBody.Password).Return("jwt", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"token\":\"jwt\"}", w.Body.String())
}

func TestSignIn_Failed_400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	body, err := json.Marshal("not json")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
func TestSignIn_ErrUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	signInBody := &models.SignInput{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(signInBody)
	assert.NoError(t, err)

	uc.On("SignIn", signInBody.Username, signInBody.Password).Return("", auth.ErrUserNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"message\":\"user not found\"}", w.Body.String())
}

func TestSignIn_ErrUnknown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	signInBody := &models.SignInput{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(signInBody)
	assert.NoError(t, err)

	uc.On("SignIn", signInBody.Username, signInBody.Password).Return("", auth.ErrUnauthorized)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"message\":\"unknown error\"}", w.Body.String())
}

func TestChangePassword_ErrInvalidCreds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	changePassBody := &models.ChangePasswordInput{
		Username:    "testuser",
		Password:    "newpass",
		OldPassword: "testpass",
	}

	body, err := json.Marshal(changePassBody)
	assert.NoError(t, err)

	uc.On("ChangePassword", changePassBody.Username, changePassBody.OldPassword, changePassBody.Password).Return(gorm.ErrInvalidTransaction)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/change-pass", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"message\":\"invalid credentials\"}", w.Body.String())
}

func TestChangePassword_ErrUnknown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	changePassBody := &models.ChangePasswordInput{
		Username:    "testuser",
		Password:    "newpass",
		OldPassword: "testpass",
	}

	body, err := json.Marshal(changePassBody)
	assert.NoError(t, err)

	uc.On("ChangePassword", changePassBody.Username, changePassBody.OldPassword, changePassBody.Password).Return(auth.ErrUnauthorized)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/change-pass", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"message\":\"unknown error\"}", w.Body.String())
}

func TestChangePassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	changePassBody := &models.ChangePasswordInput{
		Username:    "testuser",
		Password:    "newpass",
		OldPassword: "testpass",
	}

	body, err := json.Marshal(changePassBody)
	assert.NoError(t, err)

	uc.On("ChangePassword", changePassBody.Username, changePassBody.OldPassword, changePassBody.Password).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/change-pass", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"Password berhasil diubah\"}", w.Body.String())
}

func TestChangePassword_Failed_400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	body, err := json.Marshal("not json")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/change-pass", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
func TestDeleteUser_Sucess_200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	deleteMeBody := &models.DeleteInput{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(deleteMeBody)
	assert.NoError(t, err)

	uc.On("DeleteAccount", deleteMeBody.Username, deleteMeBody.Password).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/delete-me", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"Akun berhasil dihapus\"}", w.Body.String())
}

func TestDeleteUser_Failed_400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	body, err := json.Marshal("not json")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/delete-me", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
func TestDeleteUser_ErrUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	deleteMeBody := &models.DeleteInput{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(deleteMeBody)
	assert.NoError(t, err)

	uc.On("DeleteAccount", deleteMeBody.Username, deleteMeBody.Password).Return(errors.New("record not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/delete-me", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"message\":\"user not found\"}", w.Body.String())
}

func TestDeleteUser_ErrUnknown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	uc := new(mock.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	deleteMeBody := &models.DeleteInput{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(deleteMeBody)
	assert.NoError(t, err)

	uc.On("DeleteAccount", deleteMeBody.Username, deleteMeBody.Password).Return(auth.ErrUnknown)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/delete-me", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"message\":\"unknown error\"}", w.Body.String())
}
