package mock

import (
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) SignUp(inp models.SignUpInput) error {
	args := m.Called(inp.Username, inp.Email, inp.Password)

	return args.Error(0)
}

func (m *AuthUseCaseMock) SignIn(inp models.SignInput) (string, error) {
	args := m.Called(inp.Username, inp.Password)

	return args.Get(0).(string), args.Error(1)
}

func (m *AuthUseCaseMock) DeleteAccount(inp models.DeleteInput) error {
	args := m.Called(inp.Username, inp.Password)

	return args.Error(0)
}

func (m *AuthUseCaseMock) ChangePassword(inp models.ChangePasswordInput) error {
	args := m.Called(inp.Username, inp.OldPassword, inp.Password)

	return args.Error(0)
}

func (m *AuthUseCaseMock) ParseToken(accessToken string) (*models.User, error) {
	args := m.Called(accessToken)

	return args.Get(0).(*models.User), args.Error(1)
}
