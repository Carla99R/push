package repository

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"push/model"
	"push/utils"
	"strconv"
)

var (
	DB  *gorm.DB
	err error
)

var (
	OperativeSystems []model.OperativeSystems
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(utils.ReadXml("01"))
		return
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), port, os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&model.OperativeSystems{}, &model.PushNotification{})

	DB.Model(OperativeSystems).Select("*").Find(&OperativeSystems)
	fmt.Println("Database connection is successful")

}
