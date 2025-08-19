package repositories

import (
	"math"

	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindById(id uint) (*models.Course, error)
	GetAllCourses(query models.CoursesQuery) ([]models.CourseWithModulesCount, int64, error)
	FindModulesByCourseID(id uint) ([]models.Module, int64, error)
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

func (r *courseRepository) FindById(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) GetAllCourses(query models.CoursesQuery) ([]models.CourseWithModulesCount, int64, error) {
	var results []models.CourseWithModulesCount
	var totalItems int64

	base := r.db.Model(&models.Course{})

	if query.Q != "" {
		search := "%" + query.Q + "%"
		base = base.Where("courses.title ILIKE ? OR courses.instructor ILIKE ? OR EXISTS (SELECT 1 FROM unnest(courses.topics) t WHERE t ILIKE ?)",
			search, search, search)
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

	db := base.Select("courses.*, COUNT(modules.id) as modules_count").
		Joins("LEFT JOIN modules ON modules.course_id = courses.id").
		Group("courses.id")

	offset := (query.Page - 1) * query.Limit
	if err := db.Limit(query.Limit).Offset(offset).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, totalItems, nil
}

func (r *courseRepository) FindModulesByCourseID(id uint) ([]models.Module, int64, error) {
	var modules []models.Module
	var count int64

	if err := r.db.Model(&models.Module{}).
		Where("course_id = ?", id).
		Count(&count).
		Find(&modules).Error; err != nil {
		return nil, 0, err
	}

	return modules, count, nil
}
