package repository

import "feed_system/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	ListAll() ([]model.User, error)
}
