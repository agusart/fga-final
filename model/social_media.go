package model

import "time"

type SocialMedia struct {
	ID             int `gorm:"primaryKey"`
	Name           string
	SocialMediaURL string
	UserID         uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	User           User
}
