package models

import "time"

type Purchase struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint    `json:"user_id" gorm:"not null;index"`
	CourseID  uint    `json:"course_id" gorm:"not null;index"`
	Amount    float64 `json:"amount" gorm:"type:numeric(10,2);not null"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Course Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
