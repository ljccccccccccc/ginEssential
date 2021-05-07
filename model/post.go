package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserID uint `json:"user_id" gorm:"not null"`
	CategoryId uint `json:"category_id" gorm:"not null"`
	Category *Category
	Title string `json:"title" gorm:"type:varchar(50);not null"`
	HeadImg string `json:"head_img"`
	Content string `json:"content" gorm:"type:text;not null"`
}

func (post *Post) BeforeCreate(tx *gorm.DB) error {
	post.ID = uuid.NewV4()
	return nil
}
