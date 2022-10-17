package model

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PhotoID   uint
	Message   string `json:"message" valid:"required~message is blank"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User  `valid:"-"`
	Photo     Photo `valid:"-"`
}
