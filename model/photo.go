package model

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"time"
)

type Photo struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `json:"title" valid:"required~title is empty"`
	Caption   string
	PhotoURL  string `json:"photo_url"  valid:"required~photo url is empty"`
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	User      User `valid:"-"`
}

func (photo *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(photo)
	if err != nil {
		return err
	}

	return err
}

func (photo *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(photo)
	if err != nil {
		return err
	}

	return err
}
