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

type CourseFormInput struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description" binding:"required"`
	Instructor     string                `form:"instructor" binding:"required"`
	Topics         pq.StringArray        `form:"topics" binding:"required"`
	Price          float64               `form:"price" binding:"required"`
	ThumbnailImage *multipart.FileHeader `form:"thumbnail_image"`
}

type SearchQuery struct {
	Q string `form:"q"`
	PaginationQuery
}

type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (q *PaginationQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 15
	}
	if q.Limit > 50 {
		q.Limit = 50
	}
}

type ModuleFormInput struct {
	Title        string                `form:"title" binding:"required"`
	Description  string                `form:"description" binding:"required"`
	PDFContent   *multipart.FileHeader `form:"pdf_content"`
	VideoContent *multipart.FileHeader `form:"video_content"`
}

type ReorderModulesRequest struct {
	ModuleOrder []struct {
		ID    string `json:"id"`
		Order int    `json:"order"`
	} `json:"module_order"`
}

type PostBalance struct {
	Increment float64 `json:"increment"`
}

type PostUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"omitempty,min=8"`
}

type CoursesPageData struct {
	Courses    []CourseCardData
	Page       int
	TotalPages int
	TotalItems int
	Pages      []int
	Limit      int
	Search     string
	User       *User
}

type CourseCardData struct {
	ID             uint
	Title          string
	Instructor     string
	Topics         pq.StringArray
	ThumbnailImage string
	Price          float64
	Purchased      bool
}

type CourseDetailPageData struct {
	Course         *Course
	Purchased      bool
	User           *User
	CourseProgress *CourseProgress
	CertificateURL *string
}

type CourseModulesPageData struct {
	Course         *Course
	User           *User
	Modules        []ModuleWithIsCompleted
	CourseProgress CourseProgress
	CurrentModule  *ModuleWithIsCompleted
	ContentType    string
}
