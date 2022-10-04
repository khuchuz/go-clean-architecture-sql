package usecase

import (
	"testing"

	"github.com/khuchuz/go-clean-architecture-sql/auth"
	"github.com/khuchuz/go-clean-architecture-sql/auth/entities"
	"github.com/khuchuz/go-clean-architecture-sql/auth/repository/mock"
	"github.com/khuchuz/go-clean-architecture-sql/models"
	"github.com/stretchr/testify/assert"
)

func Test_SignUp_Success(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("SQLIsUserExistByUsername", username).Return(false)
	repo.On("SQLIsUserExistByEmail", email).Return(false)
	repo.On("SQLCreateUser", user).Return(nil)
	err := uc.SignUp(entities.SignUpInput{Username: username, Email: email, Password: password})
	assert.NoError(t, err)
}

func Test_SignUp_Failed_DupUsername(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("SQLIsUserExistByUsername", username).Return(true)
	repo.On("SQLCreateUser", user).Return(nil)
	err := uc.SignUp(entities.SignUpInput{Username: username, Email: email, Password: password})
	assert.Error(t, err, auth.ErrUserDuplicate)
}

func Test_SignUp_Failed_DupEmail(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("SQLIsUserExistByUsername", username).Return(false)
	repo.On("SQLIsUserExistByEmail", email).Return(true)
	repo.On("SQLCreateUser", user).Return(nil)
	err := uc.SignUp(entities.SignUpInput{Username: username, Email: email, Password: password})
	assert.Error(t, err, auth.ErrEmailDuplicate)
}
func Test_SignUp_Failed_EmptyUsername(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = ""
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("SQLCreateUser", user).Return(nil)
	err := uc.SignUp(entities.SignUpInput{Username: username, Email: email, Password: password})
	assert.Error(t, err, auth.ErrDataTidakLengkap)
}

func Test_SignUp_Failed_EmptyEmail(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = ""
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("SQLCreateUser", user).Return(nil)
	err := uc.SignUp(entities.SignUpInput{Username: username, Email: email, Password: password})
	assert.Error(t, err, auth.ErrDataTidakLengkap)
}

func Test_SignUp_Failed_Password(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = ""

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("SQLCreateUser", user).Return(nil)
	err := uc.SignUp(entities.SignUpInput{Username: username, Email: email, Password: password})
	assert.Error(t, err, auth.ErrDataTidakLengkap)
}

func Test_SignIn_Success(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign In (Get Auth Token)
	repo.On("SQLGetUser", user.Username, user.Password).Return(user, nil)
	token, err := uc.SignIn(entities.SignInput{Username: username, Password: password})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func Test_SignIn_Failed(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign In (Get Auth Token)
	repo.On("SQLGetUser", user.Username, user.Password).Return(user, auth.ErrUnknown)
	token, err := uc.SignIn(entities.SignInput{Username: username, Password: password})
	assert.Error(t, err, auth.ErrUserNotFound)
	assert.Empty(t, token)
}
func Test_ParseToken_Success(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"
		password = "pass"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	repo.On("SQLGetUser", user.Username, user.Password).Return(user, nil)
	token, err := uc.SignIn(entities.SignInput{Username: username, Password: password})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	parsedUser, err := uc.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user, parsedUser)
}

func Test_ParseToken_Failed(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		email    = "usermock@gmail.com"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)
	var token = "mboh"

	// Verify token
	parsedUser, err := uc.ParseToken(token)
	assert.Error(t, err)
	assert.NotEqual(t, user, parsedUser)
}

func Test_ChangePassword_Sucess(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username     = "usermock"
		email        = "usermock@gmail.com"
		password     = "pass"
		newpass      = "newpass"
		newpasscrypt = "128bef5c354a597f9160a61f267a7f06dab7d042"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Change Password
	repo.On("SQLGetUser", user.Username, user.Password).Return(user, nil)
	repo.On("SQLUpdatePassword", user.Username, newpasscrypt).Return(nil)
	err := uc.ChangePassword(entities.ChangePasswordInput{Username: username, OldPassword: password, Password: newpass})
	assert.NoError(t, err)
}

func Test_ChangePassword_Failed_WrongOldPass(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username     = "usermock"
		email        = "usermock@gmail.com"
		password     = "pass"
		newpass      = "newpass"
		newpasscrypt = "128bef5c354a597f9160a61f267a7f06dab7d042"

		user = &models.User{
			Username: username,
			Email:    email,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Change Password
	repo.On("SQLGetUser", user.Username, user.Password).Return(user, auth.ErrUnknown)
	repo.On("SQLUpdatePassword", user.Username, newpasscrypt).Return(nil)
	err := uc.ChangePassword(entities.ChangePasswordInput{Username: username, OldPassword: password, Password: newpass})
	assert.Error(t, err, auth.ErrUserNotFound)
}

func Test_ChangePassword_Failed_EmptyField(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		password = "pass"
		newpass  = "newpass"
	)

	// Empty Username
	err := uc.ChangePassword(entities.ChangePasswordInput{Username: "", OldPassword: password, Password: newpass})
	assert.EqualError(t, err, "data tidak lengkap")

	// Empty Password
	err = uc.ChangePassword(entities.ChangePasswordInput{Username: username, OldPassword: password, Password: ""})
	assert.EqualError(t, err, "data tidak lengkap")

	// Empty OldPassword
	err = uc.ChangePassword(entities.ChangePasswordInput{Username: username, OldPassword: "", Password: newpass})
	assert.EqualError(t, err, "data tidak lengkap")

}

func Test_ChangePassword_Failed_EqualNewOld(t *testing.T) {
	repo := new(mock.UserStorageMock)
	uc := NewAuthUseCase(repo, "salt", []byte("secret"), 86400)
	var (
		username = "usermock"
		password = "pass"
	)

	// Empty OldPassword
	err := uc.ChangePassword(entities.ChangePasswordInput{Username: username, OldPassword: password, Password: password})
	assert.EqualError(t, err, "password baru tidak boleh sama dengan password lama")
}
