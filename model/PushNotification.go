package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type PushNotification struct {
	gorm.Model
	Id                uint             `gorm:"primaryKey" gorm:"autoIncrement"`
	Title             string           `gorm:"not null"`
	Message           string           `gorm:"not null"`
	DeviceTokenId     uuid.UUID        `gorm:"not null"`
	OperativeSystemId uuid.UUID        `gorm:"not null"`
	SendId            string           `gorm:"not null"`
	OperativeSystem   OperativeSystems `gorm:"foreignKey:Id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Date              string           `gorm:"not null"`
	Time              string           `gorm:"not null"`
}

func (PushNotification) TableName() string {
	return "push_notification"
}
