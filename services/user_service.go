package services

import (
	"DistributionFlex/models"
	"DistributionFlex/repositories"
	"context"
	"errors"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Login(ctx context.Context, username, password string) (*models.User, error) {
	user, err := us.repo.FindUserByUsernameAndPassword(ctx, username, password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}
