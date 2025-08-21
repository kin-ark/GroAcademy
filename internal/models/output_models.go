package models

import "time"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

type CourseResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Instructor     string    `json:"instructor"`
	Topics         []string  `json:"topics"`
	Price          float64   `json:"price"`
	ThumbnailImage *string   `json:"thumbnail_image"`
	TotalModules   int       `json:"total_modules"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ModuleResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PDFContent   string    `json:"pdf_content"`
	VideoContent string    `json:"video_content"`
	Order        int       `json:"order"`
	IsCompleted  bool      `json:"is_completed"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BuyCourseResponse struct {
	CourseID      uint    `json:"course_id"`
	UserBalance   float64 `json:"user_balance"`
	TransactionID uint    `json:"transaction_id"`
}

type MyCoursesResponse struct {
	Course
	ProgressPercentage float64   `json:"progress_percentage"`
	PurchasedAt        time.Time `json:"purchased_at"`
}

type MarkModuleResponse struct {
	ModuleID       uint           `json:"module_id"`
	IsCompleted    bool           `json:"is_completed"`
	CourseProgress CourseProgress `json:"course_progress"`
	CertificateURL *string        `json:"certificate_url"`
}

type UsersResponse struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}
