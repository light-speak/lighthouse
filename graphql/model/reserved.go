package model

import (
	"time"

	"gorm.io/gorm"
)

type PaginateInfo struct {
	CurrentPage int  `json:"current_page"`
	TotalPage   int  `json:"total_page"`
	TotalCount  int  `json:"total_count"`
	HasNextPage bool `json:"has_next_page"`
}

type Model struct {
	ID        int64     `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModelSoftDelete struct {
	Model
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
