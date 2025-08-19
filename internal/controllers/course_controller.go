package controllers

import (
	"net/http"

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

func (courseController *CourseController) PostCourse(c *gin.Context) {
	var input models.PostCourseFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	result, err := courseController.service.CreateCourse(c, input)

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
