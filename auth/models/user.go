package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
}

type SignInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordInput struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldpassword"`
	Password    string `json:"password"`
}

type SignUpInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignResponse struct {
	Message string `json:"message"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
