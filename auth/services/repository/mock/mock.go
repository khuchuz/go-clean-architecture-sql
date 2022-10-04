package mock

import (
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"github.com/stretchr/testify/mock"
)

type UserStorageMock struct {
	mock.Mock
}

func (s *UserStorageMock) SQLCreateUser(user *models.User) error {
	args := s.Called(user)

	return args.Error(0)
}

func (s *UserStorageMock) SQLGetUser(username, password string) (*models.User, error) {
	args := s.Called(username, password)

	return args.Get(0).(*models.User), args.Error(1)
}

func (s *UserStorageMock) SQLUpdatePassword(username, oldpassword, password string) error {
	args := s.Called(username, oldpassword, password)

	return args.Error(0)
}

func (s *UserStorageMock) SQLIsUserExistByUsername(username string) bool {
	args := s.Called(username)

	return args.Bool(0)
}

func (s *UserStorageMock) SQLIsUserExistByEmail(email string) bool {
	args := s.Called(email)

	return args.Bool(0)
}

func (s *UserStorageMock) SQLDeleteUser(username, password string) error {
	args := s.Called(username, password)

	return args.Error(0)
}
