package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string  `json:"first_name" gorm:"size:100;not null"`
	LastName  string  `json:"last_name" gorm:"size:100;not null"`
	Username  string  `json:"username" gorm:"size:50;unique;not null"`
	Email     string  `json:"email" gorm:"size:100;unique;not null"`
	Password  string  `json:"-" gorm:"not null"`
	Role      string  `json:"role" gorm:"size:20;not null;default:'student'"`
	Balance   float64 `json:"balance" gorm:"type:decimal(10,2);default:0"`
}
