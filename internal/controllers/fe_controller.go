package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/services"
)

type FEController struct {
	us services.UserService
	cs services.CourseService
	ms services.ModuleService
}

func NewFEController(us services.UserService, cs services.CourseService, ms services.ModuleService) *FEController {
	return &FEController{us: us, cs: cs, ms: ms}
}

func (fc *FEController) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (fc *FEController) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func getUserFromContext(c *gin.Context) (*models.User, uint) {
	userVal, exists := c.Get("user")
	if !exists {
		return nil, 0
	}
	if u, ok := userVal.(models.User); ok {
		return &u, u.ID
	}
	return nil, 0
}

func (fc *FEController) GetCoursesPage(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	search := c.DefaultQuery("q", "")

	query := models.SearchQuery{Q: search, PaginationQuery: models.PaginationQuery{Page: page, Limit: limit}}

	courses, pagination, err := fc.cs.GetAllCourses(query)
	if err != nil {
		log.Printf("Failed to get all courses: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message":    "Could not retrieve courses.",
			"StatusCode": http.StatusInternalServerError})
		return
	}

	user, userID := getUserFromContext(c)

	var courseIDs []uint
	for _, course := range courses {
		courseIDs = append(courseIDs, course.Course.ID)
	}

	purchaseStatus, err := fc.cs.GetPurchaseStatusForCourses(courseIDs, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get purchase status for user %d: %v", userID, err)
		purchaseStatus = make(map[uint]bool)
	}

	courseCards := make([]models.CourseCardData, 0, len(courses))
	for _, course := range courses {
		courseCards = append(courseCards, models.CourseCardData{
			ID:             course.ID,
			Title:          course.Course.Title,
			Instructor:     course.Course.Instructor,
			Purchased:      purchaseStatus[course.Course.ID],
			Topics:         course.Topics,
			ThumbnailImage: course.ThumbnailImage,
			Price:          course.Price,
		})
	}

	var pages []int
	for i := 1; i <= pagination.TotalPages; i++ {
		pages = append(pages, i)
	}

	c.HTML(http.StatusOK, "courses.html", models.CoursesPageData{
		Courses:    courseCards,
		Page:       pagination.CurrentPage,
		TotalPages: pagination.TotalPages,
		TotalItems: pagination.TotalItems,
		Pages:      pages,
		Limit:      limit,
		Search:     search,
		User:       user,
	})
}

func (fc *FEController) GetMyCoursesPage(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	search := c.DefaultQuery("q", "")

	query := models.SearchQuery{Q: search, PaginationQuery: models.PaginationQuery{Page: page, Limit: limit}}

	user, userID := getUserFromContext(c)

	courses, pagination, err := fc.cs.GetCoursesByUser(user, query)
	if err != nil {
		log.Printf("Failed to get all courses: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message":    "Could not retrieve courses.",
			"StatusCode": http.StatusInternalServerError})
		return
	}

	var courseIDs []uint
	for _, course := range courses {
		courseIDs = append(courseIDs, course.Course.ID)
	}

	purchaseStatus, err := fc.cs.GetPurchaseStatusForCourses(courseIDs, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get purchase status for user %d: %v", userID, err)
		purchaseStatus = make(map[uint]bool)
	}

	courseCards := make([]models.CourseCardData, 0, len(courses))
	for _, course := range courses {
		courseCards = append(courseCards, models.CourseCardData{
			ID:             course.ID,
			Title:          course.Course.Title,
			Instructor:     course.Course.Instructor,
			Purchased:      purchaseStatus[course.Course.ID],
			Topics:         course.Topics,
			ThumbnailImage: course.ThumbnailImage,
			Price:          course.Price,
		})
	}

	var pages []int
	for i := 1; i <= pagination.TotalPages; i++ {
		pages = append(pages, i)
	}

	c.HTML(http.StatusOK, "courses.html", models.CoursesPageData{
		Courses:    courseCards,
		Page:       pagination.CurrentPage,
		TotalPages: pagination.TotalPages,
		TotalItems: pagination.TotalItems,
		Pages:      pages,
		Limit:      limit,
		Search:     search,
		User:       user,
	})
}

func (fc *FEController) GetCourseDetailPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Invalid course ID.",
			"StatusCode": http.StatusBadRequest})
		return
	}
	courseID := uint(id)

	course, err := fc.cs.GetCourseByID(courseID)
	if err != nil {
		log.Printf("ERROR: Course with ID %d not found: %v", courseID, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Message":    "Course not found.",
			"StatusCode": http.StatusNotFound})
		return
	}

	user, userID := getUserFromContext(c)

	purchased := false
	if userID != 0 {
		hasPurchased, err := fc.cs.HasPurchasedCourse(courseID, userID)
		if err != nil {
			log.Printf("ERROR: Could not check purchase status for course %d, user %d: %v", courseID, userID, err)
		} else {
			purchased = hasPurchased
		}
	}

	courseProgress, err := fc.ms.GetCourseProgress(courseID, *user)
	if err != nil {
		log.Printf("ERROR: Cannot find course progress")
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Message":    "Cannot find course progress.",
			"StatusCode": http.StatusNotFound})
		return
	}

	certificateUrl, err := fc.ms.GetCertificateURL(courseID, user.ID)
	if err != nil {
		log.Printf("ERROR: Cannot find certificate")
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Message":    "Cannot find certificate.",
			"StatusCode": http.StatusNotFound})
		return
	}

	c.HTML(http.StatusOK, "course-detail.html", models.CourseDetailPageData{
		Course:         course,
		CourseProgress: courseProgress,
		Purchased:      purchased,
		User:           user,
		CertificateURL: certificateUrl,
	})
}

func (fc *FEController) BuyCourseFE(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Invalid course ID.",
			"StatusCode": http.StatusBadRequest})
		return
	}
	courseID := uint(id)

	user, _ := getUserFromContext(c)
	if user == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Cannot get User",
			"StatusCode": http.StatusBadRequest})
		return
	}

	_, err = fc.cs.BuyCourse(courseID, user)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message":    err.Error(),
			"StatusCode": http.StatusInternalServerError,
		})
		return
	}

	c.Redirect(http.StatusFound, "/course/"+fmt.Sprint(courseID)+"/modules")
}

func (fc *FEController) GetCourseModulesPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Invalid course ID.",
			"StatusCode": http.StatusBadRequest})
		return
	}
	courseID := uint(id)

	user, _ := getUserFromContext(c)
	if user == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Cannot get User",
			"StatusCode": http.StatusBadRequest})
		return
	}

	course, err := fc.cs.GetCourseByID(courseID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message":    err.Error(),
			"StatusCode": http.StatusInternalServerError})
	}

	modules, err := fc.ms.GetModules(*user, courseID, models.PaginationQuery{})
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    err.Error(),
			"StatusCode": http.StatusInternalServerError})
	}

	courseProgress, err := fc.ms.GetCourseProgress(courseID, *user)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    err.Error(),
			"StatusCode": http.StatusInternalServerError})
	}

	var currentModule *models.ModuleWithIsCompleted
	if moduleIDStr := c.Param("moduleId"); moduleIDStr != "" {
		moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32)
		if err == nil {
			for i := range modules {
				if modules[i].ID == uint(moduleID) {
					currentModule = &modules[i]
					break
				}
			}
		}
	}

	contentType := c.Query("type")
	if contentType == "" {
		contentType = "pdf"
	}

	c.HTML(http.StatusOK, "course-modules.html", models.CourseModulesPageData{
		Course:         course,
		User:           user,
		Modules:        modules,
		CourseProgress: *courseProgress,
		CurrentModule:  currentModule,
		ContentType:    contentType,
	})
}

func (fc *FEController) ToggleModuleCompletion(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Invalid course ID.",
			"StatusCode": http.StatusBadRequest})
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Invalid module ID.",
			"StatusCode": http.StatusBadRequest})
		return
	}

	user, _ := getUserFromContext(c)
	if user == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Message":    "Cannot get User",
			"StatusCode": http.StatusBadRequest})
		return
	}

	completedStr := c.PostForm("completed")
	completed := completedStr == "true"

	err = fc.ms.ChangeModuleCompletion(uint(moduleID), *user, completed)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message":    "Failed to update completion",
			"StatusCode": http.StatusInternalServerError})
		return
	}

	redirectURL := fmt.Sprintf("/course/%d/modules/%d", courseID, moduleID)
	if contentType := c.Query("type"); contentType != "" {
		redirectURL += "?type=" + contentType
	}

	c.Redirect(http.StatusSeeOther, redirectURL)
}
