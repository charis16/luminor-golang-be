package models

import (
	"time"
)

const TableNameCategory = "categories"

// Category mapped from table <categories>
type Category struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UUID        string    `gorm:"column:uuid;default:gen_random_uuid()" json:"uuid"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Slug        string    `gorm:"column:slug;uniqueIndex" json:"slug"`
	PhotoURL    string    `gorm:"column:photo_url" json:"photo_url"`
	IsPublished bool      `gorm:"column:is_published" json:"is_published"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Category's table name
func (*Category) TableName() string {
	return TableNameCategory
}
