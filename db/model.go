package db

import "time"
import "gorm.io/gorm"

type Model struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
