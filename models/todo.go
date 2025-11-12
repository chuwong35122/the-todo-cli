package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(255);not null"`
	CompletedAt *time.Time `gorm:"default:null"`
	TagID       *uint      `gorm:"index"`
	Tag         *TodoTag   `gorm:"foreignKey:TagID;references:ID"`
}

func (i Todo) FilterValue() string { return i.Title }
func (i Todo) GetTitle() string    { return i.Title }
func (i Todo) GetTag() string      { return i.Tag.Tag }

type TodoTag struct {
	gorm.Model
	Tag   string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Todos []Todo `gorm:"foreignKey:TagID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
