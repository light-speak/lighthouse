package db

import "time"
import "gorm.io/gorm"

type Model struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ModelSoftDelete struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
