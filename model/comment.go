package model

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PhotoID   uint   `json:"photo_id" valid:"required~photo id is blank"`
	Message   string `json:"message" valid:"required~message is blank"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	User      User  `valid:"-"`
	Photo     Photo `valid:"-"`
}

func (comment *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(comment)
	if err != nil {
		return err
	}

	return err
}

func (comment *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(comment)
	if err != nil {
		return err
	}

	return err
}
