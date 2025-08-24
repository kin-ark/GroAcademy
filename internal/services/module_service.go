package services

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/repositories"
	"github.com/kin-ark/GroAcademy/internal/utils"
)

type ModuleService interface {
	CreateModule(*gin.Context, models.ModuleFormInput, uint) (*models.Module, error)
	EditModule(*gin.Context, models.ModuleFormInput, uint) (*models.Module, error)
	DeleteModuleByID(uint) error
	GetModules(user models.User, courseID uint, q models.PaginationQuery) ([]models.ModuleWithIsCompleted, error)
	BuildModuleResponses(modules []models.ModuleWithIsCompleted) []models.ModuleResponse
	GetModuleByID(id uint, user models.User) (*models.ModuleWithIsCompleted, error)
	MarkModuleAsComplete(id uint, user models.User) (*models.MarkModuleResponse, error)
	ReorderModules(req models.ReorderModulesRequest, courseID uint) error
	GetCourseProgress(id uint, user models.User) (*models.CourseProgress, error)
	ChangeModuleCompletion(moduleID uint, user models.User, completed bool) error
	GetCertificateURL(courseID, userID uint) (*string, error)
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
	baseUrl := os.Getenv("BASE_URL")
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
		module.PDFContent = baseUrl + pdfPath
	}

	if videoPath != "" {
		module.VideoContent = baseUrl + videoPath
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
		baseUrl := os.Getenv("BASE_URL")
		path := "uploads/pdf_content/" + input.PDFContent.Filename
		if err := c.SaveUploadedFile(input.PDFContent, path); err != nil {
			return nil, err
		}
		existing.PDFContent = baseUrl + path
	} else {
		existing.PDFContent = ""
	}

	if input.VideoContent != nil {
		baseUrl := os.Getenv("BASE_URL")
		path := "uploads/video_content/" + input.VideoContent.Filename
		if err := c.SaveUploadedFile(input.VideoContent, path); err != nil {
			return nil, err
		}
		existing.VideoContent = baseUrl + path
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

func (s *moduleService) GetModuleByID(id uint, user models.User) (*models.ModuleWithIsCompleted, error) {
	module, err := s.moduleRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	hasPurchased, err := s.courseRepo.HasPurchasedCourse(module.CourseID, user.ID)
	if err != nil {
		return nil, err
	}

	if !hasPurchased {
		if user.Role != "admin" {
			return nil, errors.New(user.Username + " has not bought this course!")
		} else {
			res := models.ModuleWithIsCompleted{
				Module:      *module,
				IsCompleted: false,
			}
			return &res, nil
		}
	} else {
		isCompleted, err := s.moduleRepo.IsModuleCompleted(id, user.ID)
		if err != nil {
			return nil, err
		}
		res := models.ModuleWithIsCompleted{
			Module:      *module,
			IsCompleted: isCompleted,
		}
		return &res, nil
	}
}

func (s *moduleService) MarkModuleAsComplete(id uint, user models.User) (*models.MarkModuleResponse, error) {
	err := s.moduleRepo.ChangeModuleCompletion(id, user.ID, true)
	if err != nil {
		return nil, err
	}

	module, err := s.moduleRepo.FindById(id)
	if err != nil {
		s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	isCompleted, err := s.moduleRepo.IsModuleCompleted(id, user.ID)
	if err != nil {
		s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}
	courseId := module.CourseID
	courseProgress, err := s.courseRepo.GetCourseProgress(courseId, user)
	if err != nil {
		s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	certificateURL, err := s.generateCertificateIfEligible(id, user, courseId, courseProgress)
	if err != nil {
		return nil, err
	}

	res := models.MarkModuleResponse{
		ModuleID:       module.ID,
		IsCompleted:    isCompleted,
		CourseProgress: *courseProgress,
		CertificateURL: certificateURL,
	}
	return &res, nil
}

func (s *moduleService) ReorderModules(req models.ReorderModulesRequest, courseID uint) error {
	if len(req.ModuleOrder) == 0 {
		return errors.New("module_order cannot be empty")
	}

	var moduleIDs []uint
	if err := s.moduleRepo.GetModuleIDsByCourse(courseID, &moduleIDs); err != nil {
		return errors.New("failed to fetch modules: " + err.Error())
	}

	if len(moduleIDs) == 0 {
		return errors.New("no modules found for course")
	}

	validModuleMap := make(map[uint]bool)
	for _, id := range moduleIDs {
		validModuleMap[id] = true
	}

	orderSet := make(map[int]bool)
	parsedOrders := make([]models.ModuleOrder, 0, len(req.ModuleOrder))

	for _, m := range req.ModuleOrder {
		idUint, err := strconv.ParseUint(m.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid module id format: %s", m.ID)
		}
		moduleID := uint(idUint)

		if !validModuleMap[moduleID] {
			return fmt.Errorf("invalid module id: %d", moduleID)
		}
		if m.Order < 1 || m.Order > len(moduleIDs) {
			return fmt.Errorf("invalid order %d for module %d", m.Order, moduleID)
		}
		if orderSet[m.Order] {
			return fmt.Errorf("duplicate order: %d", m.Order)
		}
		orderSet[m.Order] = true

		parsedOrders = append(parsedOrders, models.ModuleOrder{
			ID:    moduleID,
			Order: m.Order,
		})
	}

	if len(req.ModuleOrder) != len(moduleIDs) {
		return fmt.Errorf("all modules must be included in the reorder request")
	}

	return s.moduleRepo.ReorderModules(courseID, parsedOrders)
}

func (s *moduleService) GetCourseProgress(id uint, user models.User) (*models.CourseProgress, error) {
	return s.courseRepo.GetCourseProgress(id, user)
}

func (s *moduleService) ChangeModuleCompletion(moduleID uint, user models.User, completed bool) error {
	if err := s.moduleRepo.ChangeModuleCompletion(moduleID, user.ID, completed); err != nil {
		return err
	}

	module, err := s.moduleRepo.FindById(moduleID)
	if err != nil {
		return err
	}
	courseId := module.CourseID

	courseProgress, err := s.courseRepo.GetCourseProgress(courseId, user)
	if err != nil {
		return err
	}

	_, err = s.generateCertificateIfEligible(moduleID, user, courseId, courseProgress)
	if err != nil {
		return err
	}

	return nil
}

func (s *moduleService) generateCertificateIfEligible(id uint, user models.User, courseId uint, courseProgress *models.CourseProgress) (*string, error) {
	if int64(courseProgress.TotalModules) == 0 || courseProgress.CompletedModules != courseProgress.TotalModules {
		return nil, nil
	}

	course, err := s.courseRepo.FindById(courseId)
	if err != nil {
		_ = s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	cert, err := s.courseRepo.FindCourseCertificate(user.ID, courseId)
	if err == nil && cert != nil {
		return &cert.FileURL, nil
	}

	img, err := utils.GenerateCertificate(
		user.Username,
		course.Title,
		course.Instructor,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		_ = s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	fileName := fmt.Sprintf("cert_user%d_course%d.png", user.ID, courseId)
	filePath := filepath.Join("uploads/certificates", fileName)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		_ = s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	f, err := os.Create(filePath)
	if err != nil {
		_ = s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		_ = s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	publicURL := os.Getenv("BASE_URL") + "uploads/certificates/" + fileName

	certificate := models.Certificate{
		UserID:   user.ID,
		CourseID: courseId,
		FileURL:  publicURL,
	}
	if err := s.courseRepo.CreateCourseCertificate(&certificate); err != nil {
		_ = s.moduleRepo.ChangeModuleCompletion(id, user.ID, false)
		return nil, err
	}

	return &certificate.FileURL, nil
}

func (s *moduleService) GetCertificateURL(courseID, userID uint) (*string, error) {
	cert, err := s.courseRepo.FindCourseCertificate(userID, courseID)
	if err != nil {
		return nil, err
	}

	if cert != nil {
		return &cert.FileURL, nil
	}

	return nil, nil
}
