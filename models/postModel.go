package models

import "time"

type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
