package models

import "time"

type Certificate struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint   `json:"user_id" gorm:"not null;index"`
	CourseID  uint   `json:"course_id" gorm:"not null;index"`
	FileURL   string `json:"file_url" gorm:"size:255;not null"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Course    Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
