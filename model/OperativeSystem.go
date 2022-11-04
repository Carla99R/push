package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type OperativeSystems struct {
	gorm.Model
	Id              uuid.UUID `gorm:"primaryKey" gorm:"autoIncrement"`
	OperativeSystem string    `gorm:"unique" gorm:"not null"`
	IsAvailable     bool      `gorm:"default: true"`
}

func (OperativeSystems) TableName() string {
	return "operative_systems"
}
