package model

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"time"
)

type SocialMedia struct {
	ID             uint      `gorm:"primaryKey"`
	Name           string    `json:"name"  valid:"required~name is empty"`
	SocialMediaURL string    `json:"social_media_url"  valid:"required~social media url is empty"`
	UserID         uint      `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt
	User           User `json:"User" valid:"-"`
}

func (socialMedia *SocialMedia) BeforeUpdate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(socialMedia)
	if err != nil {
		return err
	}

	return err
}

func (socialMedia *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(socialMedia)
	if err != nil {
		return err
	}

	return err
}
