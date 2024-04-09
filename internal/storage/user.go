package storage

import (
	"app/internal/storage/models"
	"errors"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserExists             = errors.New("user exists")
	ErrUserDescriptionIsEmpty = errors.New("user description is empty")
)

type UserAuth interface {
	GetUser(username string) (*models.User, error)
	CreateUser(*models.User) (int64, error)
}

type UserProfile interface {
	GetUserProfile(username string) (*models.UserProfile, error)
	UpdateUserProfile(models.UserProfile) error
}
