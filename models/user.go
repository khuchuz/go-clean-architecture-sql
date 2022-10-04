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
