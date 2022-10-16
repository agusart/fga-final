package model

import "time"

type Photo struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `json:"title" valid:"required~title is empty"`
	Caption   string
	PhotoURL  string `json:"photo_url"  valid:"required~photo url is empty"`
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
