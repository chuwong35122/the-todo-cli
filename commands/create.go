package commands

import (
	"todo/models"

	"gorm.io/gorm"
)

func createTag(db *gorm.DB, tag string) (*models.TodoTag, error) {
	t := models.TodoTag{Tag: tag}

	err := db.Save(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func create(db *gorm.DB, desc, tag string) error {
	if tag == "" {
		tag = "-"
	}

	tagModel, err := createTag(db, tag)
	if err != nil {
		return err
	}

	todo := models.Todo{
		Title: desc,
		TagID: &tagModel.ID,
	}

	return db.Create(&todo).Error
}
