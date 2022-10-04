package usecase

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/khuchuz/go-clean-architecture-sql/auth"
	"github.com/khuchuz/go-clean-architecture-sql/auth/models"
	"github.com/khuchuz/go-clean-architecture-sql/auth/services"
	"github.com/khuchuz/go-clean-architecture-sql/auth/utils"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

type AuthUseCase struct {
	userRepo       services.UserRepositorySQL
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(
	userRepo services.UserRepositorySQL,
	hashSalt string,
	signingKey []byte,
	tokenTTL time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTL,
	}
}

func (a *AuthUseCase) SignUp(inp models.SignUpInput) error {

	if inp.Username == "" || inp.Email == "" || inp.Password == "" {
		return auth.ErrDataTidakLengkap
	}

	if a.userRepo.SQLIsUserExistByUsername(inp.Username) {
		return auth.ErrUserDuplicate
	}

	if a.userRepo.SQLIsUserExistByEmail(inp.Email) {
		return auth.ErrEmailDuplicate
	}

	user := &models.User{
		Username: inp.Username,
		Email:    inp.Email,
		Password: utils.HashThis(inp.Password, a.hashSalt),
	}

	return a.userRepo.SQLCreateUser(user)
}

func (a *AuthUseCase) SignIn(inp models.SignInput) (string, error) {

	password := utils.HashThis(inp.Password, a.hashSalt)

	user, err := a.userRepo.SQLGetUser(inp.Username, password)
	if err != nil {
		return "", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.signingKey)
}

func (a *AuthUseCase) ChangePassword(inp models.ChangePasswordInput) error {
	if inp.Username == "" || inp.OldPassword == "" || inp.Password == "" {
		return auth.ErrDataTidakLengkap
	}
	if inp.OldPassword == inp.Password {
		return auth.ErrPasswordSame
	}

	oldpassword := utils.HashThis(inp.OldPassword, a.hashSalt)
	password := utils.HashThis(inp.Password, a.hashSalt)

	return a.userRepo.SQLUpdatePassword(inp.Username, oldpassword, password)
}

func (a *AuthUseCase) DeleteAccount(inp models.DeleteInput) error {
	password := utils.HashThis(inp.Password, a.hashSalt)

	err := a.userRepo.SQLDeleteUser(inp.Username, password)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) ParseToken(accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
