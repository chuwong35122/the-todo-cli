package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title          string     `gorm:"column:title;type:varchar(255);not null"`
	CompletedAt    *time.Time `gorm:"column:completed_at;default:null"`
	TagID          *uint      `gorm:"column:tag_id;index"`
	Tag            *TodoTag   `gorm:"foreignKey:TagID;references:ID"`
	LastCreated    bool       `gorm:"column:last_created"`
	DisplayedTitle string     `gorm:"column:displayed_title"`
}

func (i Todo) FilterValue() string { return i.Title }

func (i Todo) GetTitle() string { return i.Title }

func (i Todo) GetTag() string {
	if i.Tag == nil {
		return "—"
	}
	return i.Tag.Tag
}

func (i Todo) GetDisplayTitle(isNew bool) string {
	if !isNew {
		return i.Title
	}
	return fmt.Sprintf("%s (new)", i.Title)
}

type TodoTag struct {
	gorm.Model
	Tag   string `gorm:"column:tag;type:varchar(100);uniqueIndex;not null"`
	Todos []Todo `gorm:"foreignKey:TagID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
