package models

import "time"

type ModuleProgress struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uint   `json:"user_id" gorm:"not null;index:idx_user_module,unique"`
	ModuleID    uint   `json:"module_id" gorm:"not null;index:idx_user_module,unique"`
	IsCompleted bool   `json:"is_completed" gorm:"default:false"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Module      Module `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ReorderModulesResponse struct {
	ModuleOrder []struct {
		ID    uint `json:"id"`
		Order int  `json:"order"`
	} `json:"module_order"`
}
