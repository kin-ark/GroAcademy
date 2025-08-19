package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Title          string         `json:"title" gorm:"size:200;not null"`
	Description    string         `json:"description" gorm:"type:text;not null"`
	Instructor     string         `json:"instructor" gorm:"size:100;not null"`
	Topics         pq.StringArray `json:"topics" gorm:"type:text[];not null;default:'{}'"`
	Price          float64        `json:"price" gorm:"type:numeric(10,2);not null"`
	ThumbnailImage string         `json:"thumbnail_image" gorm:"size:255"`
}

type CourseWithModulesCount struct {
	Course
	ModulesCount int64 `json:"modules_count"`
}
