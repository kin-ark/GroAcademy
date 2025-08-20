package services

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
)

type ModuleService interface {
	CreateModule(*gin.Context, models.ModuleFormInput, uint) (*models.Module, error)
	EditModule(*gin.Context, models.ModuleFormInput, uint) (*models.Module, error)
	DeleteModuleByID(uint) error
	GetModules(user models.User, courseID uint, q models.PaginationQuery) ([]models.ModuleWithIsCompleted, error)
	BuildModuleResponses(modules []models.ModuleWithIsCompleted) []models.ModuleResponse
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

func (s *moduleService) EditModule(c *gin.Context, input models.ModuleFormInput, id uint) (*models.Module, error) {
	existing, err := s.moduleRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if existing.PDFContent != "" {
		_ = os.Remove(existing.PDFContent)
	}

	if existing.VideoContent != "" {
		_ = os.Remove(existing.VideoContent)
	}

	if input.PDFContent != nil {
		path := "uploads/thumbnails/" + input.PDFContent.Filename
		if err := c.SaveUploadedFile(input.PDFContent, path); err != nil {
			return nil, err
		}
		existing.PDFContent = path
	} else {
		existing.PDFContent = ""
	}

	if input.VideoContent != nil {
		path := "uploads/thumbnails/" + input.VideoContent.Filename
		if err := c.SaveUploadedFile(input.VideoContent, path); err != nil {
			return nil, err
		}
		existing.VideoContent = path
	} else {
		existing.VideoContent = ""
	}

	existing.Title = input.Title
	existing.Description = input.Description

	if err := s.moduleRepo.Update(existing); err != nil {
		return nil, err
	}

	updated, err := s.moduleRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *moduleService) DeleteModuleByID(id uint) error {
	existing, err := s.moduleRepo.FindById(id)
	if err != nil {
		return err
	}

	pdfContentPath := existing.PDFContent
	videoContentPath := existing.VideoContent

	if err := s.moduleRepo.Delete(existing); err != nil {
		return err
	}

	if pdfContentPath != "" {
		_ = os.Remove(existing.PDFContent)
	}

	if videoContentPath != "" {
		_ = os.Remove(existing.VideoContent)
	}

	return nil
}

func (s *moduleService) GetModules(user models.User, courseID uint, q models.PaginationQuery) ([]models.ModuleWithIsCompleted, error) {
	var res []models.ModuleWithIsCompleted

	hasPurchased, err := s.courseRepo.HasPurchasedCourse(courseID, user.ID)
	if err != nil {
		return nil, err
	}

	if !hasPurchased {
		if user.Role != "admin" {
			return nil, errors.New(user.Username + " has not bought this course!")
		} else {
			modules, _, err := s.courseRepo.FindModulesByCourseID(courseID)
			if err != nil {
				return nil, err
			}

			for _, module := range modules {
				res = append(res, models.ModuleWithIsCompleted{
					Module:      module,
					IsCompleted: false,
				})
			}

			return res, nil
		}
	} else {
		res, err = s.courseRepo.FindModulesWithProgress(courseID, user.ID)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (s *moduleService) BuildModuleResponses(modules []models.ModuleWithIsCompleted) []models.ModuleResponse {
	var responses []models.ModuleResponse

	for _, m := range modules {
		responses = append(responses, models.ModuleResponse{
			ID:           m.ID,
			Title:        m.Title,
			Description:  m.Description,
			PDFContent:   m.PDFContent,
			VideoContent: m.VideoContent,
			Order:        m.Order,
			IsCompleted:  m.IsCompleted,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		})
	}

	return responses
}
