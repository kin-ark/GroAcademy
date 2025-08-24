package services

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
)

type CourseService interface {
	CreateCourse(*gin.Context, models.CourseFormInput) (*models.Course, error)
	EditCourse(c *gin.Context, id uint, input models.CourseFormInput) (*models.Course, error)
	GetAllCourses(query models.SearchQuery) ([]models.CourseWithModulesCount, models.PaginationResponse, error)
	GetCourseByID(id uint) (*models.Course, error)
	BuildCourseResponses(courses []models.CourseWithModulesCount) []models.CourseResponse
	GetModulesByCourse(id uint) ([]models.Module, int64, error)
	DeleteCourseByID(id uint) error
	BuyCourse(id uint, user *models.User) (*models.BuyCourseResponse, error)
	GetCoursesByUser(user *models.User, query models.SearchQuery) ([]models.MyCoursesResponse, models.PaginationResponse, error)
	HasPurchasedCourse(uint, uint) (bool, error)
	GetPurchaseStatusForCourses(courseIDs []uint, userID uint) (map[uint]bool, error)
}

type courseService struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(r repositories.CourseRepository) CourseService {
	return &courseService{courseRepo: r}
}

func (s *courseService) CreateCourse(c *gin.Context, input models.CourseFormInput) (*models.Course, error) {
	course := models.Course{Title: input.Title, Description: input.Description, Instructor: input.Instructor, Topics: input.Topics, Price: input.Price}
	if input.ThumbnailImage != nil {
		path := "uploads/thumbnails/" + input.ThumbnailImage.Filename
		err := c.SaveUploadedFile(input.ThumbnailImage, path)
		if err != nil {
			return nil, err
		}

		baseUrl := os.Getenv("BASE_URL")
		course.ThumbnailImage = baseUrl + path
	}

	if err := s.courseRepo.Create(&course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (s *courseService) GetAllCourses(query models.SearchQuery) ([]models.CourseWithModulesCount, models.PaginationResponse, error) {
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

		baseUrl := os.Getenv("BASE_URL")
		existing.ThumbnailImage = baseUrl + path
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

	updated, err := s.courseRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *courseService) DeleteCourseByID(id uint) error {
	existing, err := s.courseRepo.FindById(id)
	if err != nil {
		return err
	}

	thumbnailImagePath := existing.ThumbnailImage

	if err := s.courseRepo.Delete(existing); err != nil {
		return err
	}

	if thumbnailImagePath != "" {
		_ = os.Remove(existing.ThumbnailImage)
	}

	return nil
}

func (s *courseService) BuyCourse(id uint, user *models.User) (*models.BuyCourseResponse, error) {
	purchased, err := s.courseRepo.HasPurchasedCourse(id, user.ID)
	if err != nil {
		return nil, err
	}

	if purchased {
		return nil, errors.New(user.Username + "already purchased course: " + fmt.Sprint(id))
	} else {
		course, err := s.courseRepo.FindById(id)
		if err != nil {
			return nil, err
		}

		if course.Price > user.Balance {
			return nil, errors.New(user.Username + "balance is not enough to buy this course: " + fmt.Sprint(id))
		}

		transaction, err := s.courseRepo.BuyCourse(user, course)
		if err != nil {
			return nil, err
		}

		res := models.BuyCourseResponse{
			TransactionID: transaction.ID,
			CourseID:      id,
			UserBalance:   user.Balance,
		}
		return &res, nil
	}
}

func (s *courseService) GetCoursesByUser(user *models.User, query models.SearchQuery) ([]models.MyCoursesResponse, models.PaginationResponse, error) {
	query.Normalize()

	courses, totalItems, err := s.courseRepo.GetCoursesByUser(*user, query)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(query.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	if query.Page > totalPages {
		query.Page = totalPages
	}

	pagination := models.PaginationResponse{
		CurrentPage: query.Page,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
	}

	return courses, pagination, nil
}

func (s *courseService) HasPurchasedCourse(id uint, userId uint) (bool, error) {
	return s.courseRepo.HasPurchasedCourse(id, userId)
}

func (s *courseService) GetPurchaseStatusForCourses(courseIDs []uint, userID uint) (map[uint]bool, error) {
	purchasedIDs, err := s.courseRepo.FindPurchasedCourseIDs(userID, courseIDs)
	if err != nil {
		return nil, err
	}

	statusMap := make(map[uint]bool, len(courseIDs))
	for _, id := range courseIDs {
		statusMap[id] = false
	}
	for _, id := range purchasedIDs {
		statusMap[id] = true
	}

	return statusMap, nil
}
