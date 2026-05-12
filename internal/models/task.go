package models

import "time"

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Priority    string     `gorm:"default:'media'" json:"priority"` // baja, media, alta
	DueDate     *time.Time `json:"due_date"`
	ColumnID    uint       `gorm:"not null" json:"column_id"`
	Column      Column     `gorm:"foreignKey:ColumnID" json:"column,omitempty"`
	Order       int        `gorm:"default:0" json:"order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
