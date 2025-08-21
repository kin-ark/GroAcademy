package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/services"
)

type ModuleController struct {
	service services.ModuleService
}

func NewModuleController(s services.ModuleService) ModuleController {
	return ModuleController{service: s}
}

func (mc *ModuleController) PostModule(c *gin.Context) {
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

	var input models.ModuleFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	result, err := mc.service.CreateModule(c, input, uint(id))

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
			"id":            result.ID,
			"course_id":     result.CourseID,
			"title":         result.Title,
			"description":   result.Description,
			"order":         result.Order,
			"pdf_content":   result.PDFContent,
			"video_content": result.VideoContent,
			"created_at":    result.CreatedAt,
			"updated_at":    result.UpdatedAt,
		},
		"message": "Post module success",
		"status":  "success",
	})
}

func (mc *ModuleController) PutModule(c *gin.Context) {
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

	var input models.ModuleFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	result, err := mc.service.EditModule(c, input, uint(id))

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
			"id":            result.ID,
			"course_id":     result.CourseID,
			"title":         result.Title,
			"description":   result.Description,
			"order":         result.Order,
			"pdf_content":   result.PDFContent,
			"video_content": result.VideoContent,
			"created_at":    result.CreatedAt,
			"updated_at":    result.UpdatedAt,
		},
		"message": "Post module success",
		"status":  "success",
	})
}

func (mc *ModuleController) DeleteModuleByID(c *gin.Context) {
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

	res := mc.service.DeleteModuleByID(uint(id))

	if res != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": res.Error(),
			"data":    nil,
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (mc *ModuleController) GetModules(c *gin.Context) {
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

	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}
	u := user.(models.User)

	result, err := mc.service.GetModules(u, uint(id), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data := mc.service.BuildModuleResponses(result)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Get Modules success",
		"data":    data,
	})
}

func (mc *ModuleController) GetModuleById(c *gin.Context) {
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

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}
	u := user.(models.User)

	res, err := mc.service.GetModuleByID(uint(id), u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            res.ID,
		"course_id":     res.CourseID,
		"title":         res.Title,
		"description":   res.Description,
		"order":         res.Order,
		"pdf_content":   res.PDFContent,
		"video_content": res.VideoContent,
		"is_completed":  res.IsCompleted,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
	})
}

func (mc *ModuleController) MarkModuleAsComplete(c *gin.Context) {
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

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}
	u := user.(models.User)

	res, err := mc.service.MarkModuleAsComplete(uint(id), u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "request success",
		"data":    res,
	})
}

func (mc *ModuleController) ReorderModules(c *gin.Context) {
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

	var req models.ReorderModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"data":    nil,
		})
		return
	}

	if len(req.ModuleOrder) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "module_order cannot be empty",
			"data":    nil,
		})
		return
	}

	err = mc.service.ReorderModules(req, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to reorder modules",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "modules reordered successfully",
		"data":    req.ModuleOrder,
	})
}
