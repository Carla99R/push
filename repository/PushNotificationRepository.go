package repository

import (
	"log"
	"push/model"
)

func InsertPushNotification(data model.PushNotification) (code int, dataResponse model.PushNotification, er error) {
	result := DB.Create(&data)
	if result.Error != nil {
		log.Print(result.Error)
		return 500, data, result.Error
	}
	return 200, data, nil
}
