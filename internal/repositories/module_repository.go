package repositories

import (
	"log"

	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	Create(*models.Module) error
	Update(*models.Module) error
	Delete(*models.Module) error
	FindById(uint) (*models.Module, error)
	// MarkModuleAsComplete(id uint) error
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository() ModuleRepository {
	return &moduleRepository{db: database.DB}
}

func (r *moduleRepository) Create(module *models.Module) error {
	var maxOrder int
	err := r.db.Model(&models.Module{}).
		Where("course_id = ?", module.CourseID).
		Select(`COALESCE(MAX("order"), 0)`).
		Scan(&maxOrder).Error
	if err != nil {
		return err
	}
	module.Order = maxOrder + 1
	log.Println(maxOrder)
	return r.db.Create(module).Error
}

func (r *moduleRepository) Update(module *models.Module) error {
	return r.db.Model(&models.Module{}).
		Where("id = ?", module.ID).
		Select("*").
		Updates(module).Error
}

func (r *moduleRepository) Delete(module *models.Module) error {
	return r.db.Delete(module).Error
}

func (r *moduleRepository) FindById(id uint) (*models.Module, error) {
	var module models.Module
	err := r.db.First(&module, id).Error
	if err != nil {
		return nil, err
	}

	return &module, nil
}
