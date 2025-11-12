package commands

import (
	"todo/models"

	"gorm.io/gorm"
)

func read(db *gorm.DB) *[]models.Todo {
	var todos []models.Todo
	if err := db.Order("updated_at desc").Find(&todos).Error; err != nil {
		return nil
	}

	return &todos
}
