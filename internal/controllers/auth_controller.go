package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kin-ark/GroAcademy/internal/models"
	"github.com/kin-ark/GroAcademy/internal/services"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(s services.AuthService) AuthController {
	return AuthController{service: s}
}

func (authController *AuthController) Register(c *gin.Context) {
	var body models.RegisterInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "There's something wrong with the body",
			"data":    nil,
		})
		return
	}

	result, err := authController.service.RegisterUser(body)

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
			"id":         result.ID,
			"username":   result.Username,
			"first_name": result.FirstName,
			"last_name":  result.LastName,
		},
		"message": "User signed up",
		"status":  "success",
	})
}

func (authController *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"data":    nil,
		})
		return
	}

	token, username, err := authController.service.LoginUser(input)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "An internal error occured",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"username": username,
			"token":    token,
		},
	})
}

func (authController *AuthController) GetSelf(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Error: Unauthorized",
			"data":    nil,
		})
		return
	}

	u := user.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Request success",
		"data": gin.H{
			"id":         u.ID,
			"username":   u.Username,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"email":      u.Email,
			"balance":    u.Balance,
		},
	})
}
