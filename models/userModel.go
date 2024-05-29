package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool   `gorm:"default:false"`
	Verhash   string `gorm:"size:256;not null"`
}
