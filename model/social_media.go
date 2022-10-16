package model

import "time"

type SocialMedia struct {
	ID             int       `gorm:"primaryKey"`
	Name           string    `json:"name"  valid:"required~name is empty"`
	SocialMediaURL string    `json:"social_media_url"  valid:"required~social media url is empty"`
	UserID         uint      `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           User      `json:"User" valid:"-"`
}
