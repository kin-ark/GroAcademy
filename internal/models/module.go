package models

import (
	"time"
)

type Module struct {
	ID           uint   `gorm:"primaryKey"`
	CourseID     uint   `json:"course_id" gorm:"not null;index;uniqueIndex:idx_course_order"`
	Title        string `json:"title" gorm:"size:200;not null"`
	Description  string `json:"description" gorm:"type:text;not null"`
	PDFContent   string `json:"pdf_content" gorm:"size:255"`
	VideoContent string `json:"video_content" gorm:"size:255"`
	Order        int    `json:"order" gorm:"not null;uniqueIndex:idx_course_order"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Course Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ModuleWithIsCompleted struct {
	Module
	IsCompleted bool `json:"is_completed"`
}
