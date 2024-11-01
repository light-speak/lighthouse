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

type EnumInterface interface {
	ToString() string
}

type ModelInterface interface {
	IsModel() bool
	GetId() int64
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *gorm.DeletedAt
	TypeName() string
	TableName() string
}

type Model struct {
	Id        int64     `json:"id" gorm:"primary_key;autoIncrement"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Model) GetId() int64                  { return m.Id }
func (m *Model) GetCreatedAt() time.Time       { return m.CreatedAt }
func (m *Model) GetUpdatedAt() time.Time       { return m.UpdatedAt }
func (m *Model) GetDeletedAt() *gorm.DeletedAt { return nil }
func (m *Model) TypeName() string              { return "Model" }
func (m *Model) TableName() string             { return "models" }

type ModelSoftDelete struct {
	Model
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *ModelSoftDelete) GetDeletedAt() *gorm.DeletedAt { return &m.DeletedAt }
