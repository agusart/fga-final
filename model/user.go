package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username" gorm:"index:idx_username,unique" valid:"required~username is empty"`
	Password  string    `json:"password" valid:"required~password is empty,minstringlength(6)~minimum password must 6 char length"`
	Email     string    `json:"email" gorm:"index:idx_email,unique" valid:"required~email is blank,email~email is invalid"`
	Age       int       `json:"age" valid:"required~age is empty,range(8|125)~age is must between 8-125"`
	Photos    []Photo
}
