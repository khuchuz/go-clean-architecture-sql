package repository

import (
	"context"

	"github.com/khuchuz/go-clean-architecture-sql/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	return &UserRepository{
		db: db.Collection(collection),
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	model := toMongoUser(user)
	res, err := r.db.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r UserRepository) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	user := new(User)
	err := r.db.FindOne(ctx, bson.M{
		"username": username,
		"password": password,
	}).Decode(user)

	if err != nil {
		return nil, err
	}

	return toModel(user), nil
}

func (r UserRepository) UpdatePassword(ctx context.Context, username, password string) error {
	_, err := r.db.UpdateOne(ctx,
		bson.M{"username": username},
		bson.D{
			{"$set", bson.D{{"password", password}}},
		})
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepository) IsUserExistByUsername(ctx context.Context, username string) bool {
	user := new(User)
	err := r.db.FindOne(ctx, bson.M{
		"username": username,
	}).Decode(user)

	if err != nil {
		return false
	}
	return true
}

func (r UserRepository) IsUserExistByEmail(ctx context.Context, email string) bool {
	user := new(User)
	err := r.db.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(user)

	if err != nil {
		return false
	}
	return true
}

func toMongoUser(u *models.User) *User {
	return &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

func toModel(u *User) *models.User {
	return &models.User{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}
