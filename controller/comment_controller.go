package controller

import "gorm.io/gorm"

type commentController struct {
	db *gorm.DB
}

func NewCommentController(db *gorm.DB) *commentController {
	return &commentController{
		db: db,
	}
}
