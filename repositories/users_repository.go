package repositories

import (
	"DistributionFlex/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db.Collection("users"),
	}
}

func (ur *UserRepository) FindUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.User, error) {
	filter := bson.M{"username": username}
	user := &models.User{}
	err := ur.db.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("passwords do not match")
	}
	return user, nil
}

func (ur *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.M{"email": email}
	user := &models.User{}
	err := ur.db.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (ur *UserRepository) UpdatePasswordByEmail(ctx context.Context, email, newPasswordHash string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"passwordHash": newPasswordHash}}

	_, err := ur.db.UpdateOne(ctx, filter, update)
	return err
}
