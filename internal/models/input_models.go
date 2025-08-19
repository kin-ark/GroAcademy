package models

import (
	"mime/multipart"

	"github.com/lib/pq"
)

type RegisterInput struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,eqfield=Password"`
}

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type PostCourseFormInput struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description" binding:"required"`
	Instructor     string                `form:"instructor" binding:"required"`
	Topics         pq.StringArray        `form:"topics" binding:"required"`
	Price          float64               `form:"price" binding:"required"`
	ThumbnailImage *multipart.FileHeader `form:"thumbnail_image"`
}
