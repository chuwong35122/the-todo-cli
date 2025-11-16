package commands

import (
	"todo/models"

	"gorm.io/gorm"
)

func read(db *gorm.DB, limit, offset int, completed bool) (*[]models.Todo, error) {
	var todos []models.Todo

	query := db.Preload("Tag").
		Order("updated_at DESC")

	if !completed {
		query = query.Where("completed_at IS NULL")
	} else {
		query = query.Where("completed_at IS NOT NULL")
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&todos).Error; err != nil {
		return nil, err
	}

	return &todos, nil
}

func readAll(db *gorm.DB, limit, offset int) (*[]models.Todo, error) {
	var todos []models.Todo

	if err := db.Preload("Tag").
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&todos).Error; err != nil {
		return nil, err
	}

	return &todos, nil
}

func countAll(db *gorm.DB) (int64, error) {
	var count int64

	if err := db.Model(&models.Todo{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
