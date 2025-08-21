package repositories

import (
	"errors"
	"fmt"
	"math"

	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	FindByIdentifier(identifier string) (*models.User, error)
	GetAllUsers(query models.SearchQuery) ([]models.User, int64, error)
	FindById(id uint) (*models.User, error)
	GetNumberOfCoursePurchased(id uint) (int, error)
	AddUserBalance(id uint, increment float64) error
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

func (r *userRepository) GetAllUsers(query models.SearchQuery) ([]models.User, int64, error) {
	var results []models.User
	var totalItems int64

	base := r.db.Model(&models.User{})

	if query.Q != "" {
		search := "%" + query.Q + "%"
		base = base.Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR users.username ILIKE ? OR users.email ILIKE ?",
			search, search, search, search)
	}

	if err := base.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	if query.Page > totalPages {
		query.Page = totalPages
	}
	if query.Page < 1 {
		query.Page = 1
	}

	offset := (query.Page - 1) * query.Limit
	if err := base.Limit(query.Limit).Offset(offset).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, totalItems, nil
}

func (r *userRepository) FindById(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetNumberOfCoursePurchased(id uint) (int, error) {
	var coursesPurchased int64
	err := r.db.Model(&models.Purchase{}).
		Where("purchases.user_id = ?", id).
		Count(&coursesPurchased).Error
	if err != nil {
		return -1, err
	}

	return int(coursesPurchased), nil

}

func (r *userRepository) AddUserBalance(id uint, increment float64) error {
	result := r.db.Debug().
		Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", increment))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no user record found for user " + fmt.Sprint(id))
	}
	return nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Select("*").
		Updates(user).Error
}
