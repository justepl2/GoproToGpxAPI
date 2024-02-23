package db

import (
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Save(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}
