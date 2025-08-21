package repositories

import (
	"errors"
	"fmt"
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
	IsModuleCompleted(id uint, userId uint) (bool, error)
	MarkModuleAsComplete(moduleID uint, userID uint) error
	ReorderModules(courseID uint, orders []models.ModuleOrder) error
	GetModuleIDsByCourse(courseID uint, ids *[]uint) error
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
	courseID := module.CourseID
	deletedOrder := module.Order
	err := r.db.Delete(module).Error
	if err != nil {
		return err
	}

	err = database.DB.Model(&models.Module{}).
		Where("course_id = ?", courseID).
		Where("\"order\" > ?", deletedOrder).
		Update("\"order\"", gorm.Expr("\"order\" - ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *moduleRepository) FindById(id uint) (*models.Module, error) {
	var module models.Module
	err := r.db.First(&module, id).Error
	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) IsModuleCompleted(id uint, userId uint) (bool, error) {
	var isCompleted bool
	err := r.db.Model(&models.Module{}).
		Select("module_progresses.is_completed").
		Joins("JOIN module_progresses ON modules.id = module_progresses.module_id").
		Where("modules.id = ?", id).Where("module_progresses.user_id = ?", userId).
		Scan(&isCompleted).Error
	if err != nil {
		return false, err
	}

	return isCompleted, nil
}

func (r *moduleRepository) MarkModuleAsComplete(moduleID uint, userID uint) error {
	result := r.db.Model(&models.ModuleProgress{}).
		Where("module_id = ? AND user_id = ?", moduleID, userID).
		Update("is_completed", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no progress record found for user " + fmt.Sprint(userID) + "and module " + fmt.Sprint(moduleID))
	}
	return nil
}

func (r *moduleRepository) ReorderModules(courseID uint, orders []models.ModuleOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, m := range orders {
			if err := tx.Model(&models.Module{}).
				Where("id = ? AND course_id = ?", m.ID, courseID).
				Update("order", -int(m.ID)).Error; err != nil {
				return err
			}
		}

		for _, m := range orders {
			if err := tx.Model(&models.Module{}).
				Where("id = ? AND course_id = ?", m.ID, courseID).
				Update("order", m.Order).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *moduleRepository) GetModuleIDsByCourse(courseID uint, ids *[]uint) error {
	return r.db.Model(&models.Module{}).
		Where("course_id = ?", courseID).
		Pluck("id", ids).Error
}
