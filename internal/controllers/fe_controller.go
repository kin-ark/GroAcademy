package controllers

import (
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
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Message": "Could not retrieve courses."})
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

func (fc *FEController) GetCourseDetailPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"Message": "Invalid course ID."})
		return
	}
	courseID := uint(id)

	course, err := fc.cs.GetCourseByID(courseID)
	if err != nil {
		log.Printf("ERROR: Course with ID %d not found: %v", courseID, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Message": "Course not found."})
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

	c.HTML(http.StatusOK, "course-detail.html", models.CourseDetailPageData{
		Course:    course,
		Purchased: purchased,
		User:      user,
	})
}
