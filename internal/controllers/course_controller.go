package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/services"
)

type CourseController struct {
	service services.CourseService
}

func NewCourseController(s services.CourseService) CourseController {
	return CourseController{service: s}
}

func (cc *CourseController) PostCourse(c *gin.Context) {
	var input models.CourseFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	result, err := cc.service.CreateCourse(c, input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":              result.ID,
			"title":           result.Title,
			"description":     result.Description,
			"instructor":      result.Instructor,
			"topics":          result.Topics,
			"price":           result.Price,
			"thumbnail_image": result.ThumbnailImage,
			"created_at":      result.CreatedAt,
			"updated_at":      result.UpdatedAt,
		},
		"message": "Post course success",
		"status":  "success",
	})
}

func (cc *CourseController) GetAllCourses(c *gin.Context) {
	var query models.CoursesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	courses, pagination, err := cc.service.GetAllCourses(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	coursesResponse := cc.service.BuildCourseResponses(courses)

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Successfully get all courses",
		"data":       coursesResponse,
		"pagination": pagination,
	})
}

func (cc *CourseController) GetCourseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid course ID",
			"data":    nil,
		})
		return
	}

	course, err := cc.service.GetCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Course not found",
			"data":    nil,
		})
		return
	}
	_, moduleCount, err := cc.service.GetModulesByCourse(course.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch modules count",
			"data":    nil,
		})
		return
	}

	courseResponse := models.CourseResponse{
		ID:             course.ID,
		Title:          course.Title,
		Description:    course.Description,
		Instructor:     course.Instructor,
		Topics:         course.Topics,
		Price:          course.Price,
		ThumbnailImage: &course.ThumbnailImage,
		TotalModules:   int(moduleCount),
		CreatedAt:      course.CreatedAt,
		UpdatedAt:      course.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully retrieved course",
		"data":    courseResponse,
	})
}

func (cc *CourseController) PutCourse(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid course ID",
			"data":    nil,
		})
		return
	}

	var input models.CourseFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	result, err := cc.service.EditCourse(c, uint(id), input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":              result.ID,
			"title":           result.Title,
			"description":     result.Description,
			"instructor":      result.Instructor,
			"topics":          result.Topics,
			"price":           result.Price,
			"thumbnail_image": result.ThumbnailImage,
			"created_at":      result.CreatedAt,
			"updated_at":      result.UpdatedAt,
		},
		"message": "Post course success",
		"status":  "success",
	})
}
