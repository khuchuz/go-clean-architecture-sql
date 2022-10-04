package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type Register struct {
	Username string `bson:"username"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
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
