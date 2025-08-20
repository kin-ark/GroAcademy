package services

import (
	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
)

type ModuleService interface {
	CreateModule(*gin.Context, models.ModuleFormInput, uint) (*models.Module, error)
}

type moduleService struct {
	moduleRepo repositories.ModuleRepository
	courseRepo repositories.CourseRepository
}

func NewModuleService(mr repositories.ModuleRepository, cr repositories.CourseRepository) ModuleService {
	return &moduleService{moduleRepo: mr, courseRepo: cr}
}

func (s *moduleService) CreateModule(c *gin.Context, input models.ModuleFormInput, courseId uint) (*models.Module, error) {
	_, err := s.courseRepo.FindById(courseId)
	if err != nil {
		return nil, err
	}

	var pdfPath string
	var videoPath string
	if input.PDFContent != nil {
		pdfPath = "uploads/pdf_content/" + input.PDFContent.Filename
		err := c.SaveUploadedFile(input.PDFContent, pdfPath)
		if err != nil {
			return nil, err
		}
	}
	if input.VideoContent != nil {
		videoPath = "uploads/video_content/" + input.VideoContent.Filename
		err := c.SaveUploadedFile(input.VideoContent, videoPath)
		if err != nil {
			return nil, err
		}
	}

	module := models.Module{
		Title:       input.Title,
		Description: input.Description,
	}

	if pdfPath != "" {
		module.PDFContent = pdfPath
	}

	if videoPath != "" {
		module.VideoContent = videoPath
	}

	module.CourseID = courseId

	if err := s.moduleRepo.Create(&module); err != nil {
		return nil, err
	}

	return &module, nil
}
