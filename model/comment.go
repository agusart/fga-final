package model

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PhotoID   uint
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
