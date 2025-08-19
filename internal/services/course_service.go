package services

import (
	"math"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
)

type CourseService interface {
	CreateCourse(*gin.Context, models.CourseFormInput) (*models.Course, error)
	EditCourse(c *gin.Context, id uint, input models.CourseFormInput) (*models.Course, error)
	GetAllCourses(query models.CoursesQuery) ([]models.CourseWithModulesCount, models.PaginationResponse, error)
	GetCourseByID(id uint) (*models.Course, error)
	BuildCourseResponses(courses []models.CourseWithModulesCount) []models.CourseResponse
	GetModulesByCourse(id uint) ([]models.Module, int64, error)
}

type courseService struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(r repositories.CourseRepository) CourseService {
	return &courseService{courseRepo: r}
}

func (s *courseService) CreateCourse(c *gin.Context, input models.CourseFormInput) (*models.Course, error) {
	var course models.Course
	if input.ThumbnailImage != nil {
		path := "uploads/thumbnails/" + input.ThumbnailImage.Filename
		err := c.SaveUploadedFile(input.ThumbnailImage, path)
		if err != nil {
			return nil, err
		}
		course = models.Course{Title: input.Title, Description: input.Description, Instructor: input.Instructor, Topics: input.Topics, Price: input.Price, ThumbnailImage: path}
	}

	if err := s.courseRepo.Create(&course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (s *courseService) GetAllCourses(query models.CoursesQuery) ([]models.CourseWithModulesCount, models.PaginationResponse, error) {
	query.Normalize()

	courses, totalItems, err := s.courseRepo.GetAllCourses(query)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))
	if query.Page > totalPages && totalPages > 0 {
		query.Page = totalPages
	}

	pagination := models.PaginationResponse{
		CurrentPage: query.Page,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
	}

	return courses, pagination, nil
}

func (s *courseService) BuildCourseResponses(courses []models.CourseWithModulesCount) []models.CourseResponse {
	var responses []models.CourseResponse

	for _, c := range courses {
		responses = append(responses, models.CourseResponse{
			ID:             c.Course.ID,
			Title:          c.Course.Title,
			Description:    c.Course.Description,
			Instructor:     c.Course.Instructor,
			Topics:         c.Course.Topics,
			Price:          c.Course.Price,
			ThumbnailImage: &c.Course.ThumbnailImage,
			TotalModules:   int(c.ModulesCount),
			CreatedAt:      c.Course.CreatedAt,
			UpdatedAt:      c.Course.UpdatedAt,
		})
	}

	return responses
}

func (s *courseService) GetCourseByID(id uint) (*models.Course, error) {
	return s.courseRepo.FindById(id)
}

func (s *courseService) GetModulesByCourse(id uint) ([]models.Module, int64, error) {
	return s.courseRepo.FindModulesByCourseID(id)
}

func (s *courseService) EditCourse(c *gin.Context, id uint, input models.CourseFormInput) (*models.Course, error) {
	existing, err := s.courseRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if existing.ThumbnailImage != "" {
		_ = os.Remove(existing.ThumbnailImage)
	}

	if input.ThumbnailImage != nil {
		path := "uploads/thumbnails/" + input.ThumbnailImage.Filename
		if err := c.SaveUploadedFile(input.ThumbnailImage, path); err != nil {
			return nil, err
		}
		existing.ThumbnailImage = path
	} else {
		existing.ThumbnailImage = ""
	}

	existing.Title = input.Title
	existing.Description = input.Description
	existing.Instructor = input.Instructor
	existing.Topics = input.Topics
	existing.Price = input.Price

	if err := s.courseRepo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}
