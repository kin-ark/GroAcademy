package repositories

import (
	"math"

	"github.com/kin-ark/GroAcademy/internal/database"
	"github.com/kin-ark/GroAcademy/internal/models"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	Update(course *models.Course) error
	Delete(course *models.Course) error
	FindById(id uint) (*models.Course, error)
	GetAllCourses(query models.CoursesQuery) ([]models.CourseWithModulesCount, int64, error)
	FindModulesByCourseID(id uint) ([]models.Module, int64, error)
	HasPurchasedCourse(courseId uint, userId uint) (bool, error)
	FindModulesWithProgress(courseID, userID uint) ([]models.ModuleWithIsCompleted, error)
	BuyCourse(user *models.User, course *models.Course) (*models.Purchase, error)
	GetCoursesByUser(user models.User) ([]models.MyCoursesResponse, error)
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

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Model(&models.Course{}).
		Where("id = ?", course.ID).
		Select("*").
		Updates(course).Error
}

func (r *courseRepository) Delete(course *models.Course) error {
	return r.db.Delete(course).Error
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

func (r *courseRepository) HasPurchasedCourse(courseId uint, userId uint) (bool, error) {
	var count int64

	err := r.db.Model(&models.Purchase{}).
		Where("user_id = ? AND course_id = ?", userId, courseId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *courseRepository) FindModulesWithProgress(courseID, userID uint) ([]models.ModuleWithIsCompleted, error) {
	var modules []models.ModuleWithIsCompleted

	err := r.db.Model(&models.Module{}).
		Select("modules.*, COALESCE(module_progresses.is_completed, false) AS completed").
		Joins("LEFT JOIN module_progresses ON module_progresses.module_id = modules.id AND module_progresses.user_id = ?", userID).
		Where("modules.course_id = ?", courseID).
		Order("modules.order ASC").
		Scan(&modules).Error

	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *courseRepository) BuyCourse(user *models.User, course *models.Course) (*models.Purchase, error) {
	user.Balance -= course.Price
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	purchase := models.Purchase{
		UserID:   user.ID,
		CourseID: course.ID,
		Amount:   course.Price,
	}
	if err := r.db.Create(&purchase).Error; err != nil {
		return nil, err
	}

	var modules []models.Module
	if err := r.db.Where("course_id = ?", course.ID).Find(&modules).Error; err != nil {
		return nil, err
	}

	progresses := make([]models.ModuleProgress, len(modules))
	for i, m := range modules {
		progresses[i] = models.ModuleProgress{
			UserID:      user.ID,
			ModuleID:    m.ID,
			IsCompleted: false,
		}
	}

	if len(progresses) > 0 {
		if err := r.db.Create(&progresses).Error; err != nil {
			return nil, err
		}
	}

	return &purchase, nil
}

func (r *courseRepository) GetCoursesByUser(user models.User) ([]models.MyCoursesResponse, error) {
	var courses []models.MyCoursesResponse
	base := r.db.Model(&models.Course{})
	db := base.Select(
		"courses.*",
		"purchases.created_at AS purchased_at",
		`CASE 
			WHEN COUNT(modules.id) = 0 THEN 0
			ELSE (SUM(CASE WHEN module_progresses.is_completed THEN 1 ELSE 0 END) * 100.0 / COUNT(modules.id))
		END AS progress_percentage`).
		Joins("JOIN purchases ON purchases.course_id = courses.id").
		Joins("LEFT JOIN modules ON modules.course_id = courses.id").
		Joins("LEFT JOIN module_progresses ON module_progresses.module_id = modules.id AND module_progresses.user_id = ?", user.ID).
		Where("purchases.user_id = ?", user.ID).
		Group("courses.id, purchases.created_at")

	if err := db.Scan(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}
