package services

import (
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
)

type CourseService interface {
	CreateCourse(*gin.Context, models.PostCourseFormInput) (*models.Course, error)
}

type courseService struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(r repositories.CourseRepository) CourseService {
	return &courseService{courseRepo: r}
}

func (s *courseService) CreateCourse(c *gin.Context, input models.PostCourseFormInput) (*models.Course, error) {

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
