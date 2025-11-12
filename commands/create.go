package commands

import (
	"todo/models"

	"gorm.io/gorm"
)

func create(db *gorm.DB, desc, tag string) error {
	if tag == "" {
		tag = "me"
	}

	var tagModel *models.TodoTag

	t := models.TodoTag{}
	err := db.Where("tag = ?", tag).First(&t).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			t = models.TodoTag{Tag: tag}
			if err := db.Create(&t).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	tagModel = &t
	todo := models.Todo{
		Title: desc,
		TagID: &tagModel.ID,
	}

	return db.Create(&todo).Error
}
