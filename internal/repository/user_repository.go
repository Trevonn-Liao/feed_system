package repository

import "feed_system/internal/model"

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	ListAll() ([]model.User, error)
}
