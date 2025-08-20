package models

import (
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	CourseID     uint   `json:"course_id" gorm:"not null;index"`
	Title        string `json:"title" gorm:"size:200;not null"`
	Description  string `json:"description" gorm:"type:text;not null"`
	PDFContent   string `json:"pdf_content" gorm:"size:255"`
	VideoContent string `json:"video_content" gorm:"size:255"`
	Order        int    `json:"order" gorm:"not null"`

	Course Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ModuleWithIsCompleted struct {
	Module
	IsCompleted bool `json:"is_completed"`
}
