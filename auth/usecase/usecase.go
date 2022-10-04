package usecase

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/khuchuz/go-clean-architecture-sql/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/khuchuz/go-clean-architecture-sql/auth"
	"github.com/khuchuz/go-clean-architecture-sql/auth/entities"
	itface "github.com/khuchuz/go-clean-architecture-sql/auth/itface"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

type AuthUseCase struct {
	userRepo       itface.UserRepositorySQL
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(
	userRepo itface.UserRepositorySQL,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (a *AuthUseCase) SignUp(inp entities.SignUpInput) error {
	pwd := sha1.New()
	pwd.Write([]byte(inp.Password))
	pwd.Write([]byte(a.hashSalt))
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
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return a.userRepo.SQLCreateUser(user)
}

func (a *AuthUseCase) SignIn(inp entities.SignInput) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(inp.Password))
	pwd.Write([]byte(a.hashSalt))
	password := fmt.Sprintf("%x", pwd.Sum(nil))

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

func (a *AuthUseCase) ChangePassword(inp entities.ChangePasswordInput) error {
	if inp.Username == "" || inp.OldPassword == "" || inp.Password == "" {
		return auth.ErrDataTidakLengkap
	}
	if inp.OldPassword == inp.Password {
		return auth.ErrPasswordSame
	}
	pwd := sha1.New()
	pwd.Write([]byte(inp.OldPassword))
	pwd.Write([]byte(a.hashSalt))
	oldpassword := fmt.Sprintf("%x", pwd.Sum(nil))

	pwd2 := sha1.New()
	pwd2.Write([]byte(inp.Password))
	pwd2.Write([]byte(a.hashSalt))
	password := fmt.Sprintf("%x", pwd2.Sum(nil))

	_, err := a.userRepo.SQLGetUser(inp.Username, oldpassword)
	if err != nil {
		return auth.ErrUserNotFound
	}
	return a.userRepo.SQLUpdatePassword(inp.Username, password)
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
