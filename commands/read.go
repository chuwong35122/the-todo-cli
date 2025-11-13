package commands

import (
	"fmt"
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

	if len(todos) > 0 { // newly created
		todos[0].LastCreated = true
		todos[0].Title = fmt.Sprintf("✨%s", todos[0].Title)
	}

	return &todos
}
