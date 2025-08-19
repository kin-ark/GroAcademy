package repositories

import (
	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindById(id string) (*models.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository() CourseRepository {
	return &courseRepository{db: database.DB}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) FindById(id string) (*models.Course, error) {
	var course models.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		return nil, err
	}

	return &course, nil
}
