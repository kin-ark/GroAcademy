package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Purchase) AfterCreate(tx *gorm.DB) error {
	return createModuleProgressesForUser(tx, p.UserID, p.CourseID)
}

func (m *Module) AfterCreate(tx *gorm.DB) error {
	var purchases []Purchase
	if err := tx.Where("course_id = ?", m.CourseID).Find(&purchases).Error; err != nil {
		return err
	}

	for _, p := range purchases {
		if err := createModuleProgressesForUser(tx, p.UserID, m.CourseID); err != nil {
			return err
		}
	}

	return nil
}

func createModuleProgressesForUser(tx *gorm.DB, userID, courseID uint) error {
	var modules []Module
	if err := tx.Where("course_id = ?", courseID).Find(&modules).Error; err != nil {
		return err
	}

	progresses := make([]ModuleProgress, 0, len(modules))
	for _, m := range modules {
		progresses = append(progresses, ModuleProgress{
			UserID:      userID,
			ModuleID:    m.ID,
			IsCompleted: false,
		})
	}

	if len(progresses) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(progresses, 100).Error; err != nil {
			return err
		}
	}

	return nil
}
