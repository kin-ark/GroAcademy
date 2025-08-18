package repositories

import (
	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByIdentifier(identifier string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: database.DB}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", identifier).Or("email = ?", identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
