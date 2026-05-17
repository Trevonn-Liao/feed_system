package mysql

import (
	"errors"

	"feed_system/internal/model"
	"feed_system/internal/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ListAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Order("id asc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
