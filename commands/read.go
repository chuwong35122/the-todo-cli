package commands

import (
	"todo/models"

	"gorm.io/gorm"
)

func read(db *gorm.DB) *[]models.Todo {
	var todos []models.Todo
	err := db.Preload("Tag").
		Order("updated_at desc").Limit(10).Find(&todos).Error

	if err != nil {
		return nil
	}

	return &todos
}
