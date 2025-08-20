package models

import (
	"gorm.io/gorm"
)

type ModuleProgress struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"not null;index"`
	ModuleID    uint   `json:"module_id" gorm:"not null;index"`
	IsCompleted bool   `json:"is_completed" gorm:"default:false"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Module      Module `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
