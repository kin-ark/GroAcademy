package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/services"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) UserController {
	return UserController{service: s}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	var query models.SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bad Request",
			"data":    nil,
		})
		return
	}

	users, pagination, err := uc.service.GetUsers(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	res := uc.service.BuildUsersResponse(users)

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Request success",
		"data":       res,
		"pagination": pagination,
	})
}

func (uc *UserController) GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
			"data":    nil,
		})
		return
	}

	user, coursePurchased, err := uc.service.GetUserById(uint(id))
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
		"message": "Request success",
		"data": gin.H{
			"id":                idParam,
			"first_name":        user.FirstName,
			"last_name":         user.LastName,
			"email":             user.Email,
			"username":          user.Username,
			"balance":           user.Balance,
			"courses_purchased": coursePurchased,
		},
	})
}

