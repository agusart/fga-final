package model

import (
	"fga-final/pkg/crypto"
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
	Username  string `json:"username" gorm:"index:idx_username,unique;not null" valid:"required~username is empty"`
	Password  string `json:"password" valid:"required~password is empty,minstringlength(6)~minimum password must 6 char length"`
	Email     string `json:"email" gorm:"index:idx_email,unique;not null" valid:"required~email is blank,email~email is invalid"`
	Age       int    `json:"age" valid:"required~age is empty,range(8|125)~age is must between 8-125"`
	Photos    []Photo
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}

	hashPassword, err := crypto.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashPassword
	return err
}

func (u *User) BeforeUpdate(tx gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}

	return err
}
