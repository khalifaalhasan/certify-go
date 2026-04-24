package models

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title         string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug          string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Content       string     `gorm:"type:text;not null" json:"content"`
	Excerpt       string     `gorm:"type:text" json:"excerpt"`
	CoverImageURL string     `gorm:"type:text" json:"cover_image_url"`
	Status        string     `gorm:"type:varchar(50);default:'draft'" json:"status"`
	AuthorID      uuid.UUID  `gorm:"type:uuid;not null" json:"author_id"`
	CategoryID    uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at"`

	Author   *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Blog) TableName() string {
	return "blogs"
}