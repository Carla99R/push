package controllers

import (
	"bytes"
	contx "context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/sideshow/apns2"
	payload2 "github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"pushMessage/model"
	"pushMessage/repository"
	"pushMessage/utils"
	"strings"
	"time"
)

type CreatePushNotification struct {
	KeyId         string    `json:"key_id,omitempty"`
	TeamId        string    `json:"team_id,omitempty"`
	Title         string    `json:"title" validate:"required"`
	Message       string    `json:"message" validate:"required"`
	DeviceToken   string    `json:"device_token" validate:"required"`
	Platform      string    `json:"platform" validate:"required"`
	DeviceTokenId uuid.UUID `json:"device_token_id" validate:"required"`
}

type PushMessageResponse struct {
	Status int    `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error"`
}

type IError struct {
	Field string
	Tag   string
}

var validate *validator.Validate

func ValidateStruct(structData CreatePushNotification) map[string]interface{} {
	var errors []IError
	validate = validator.New()
	err := validate.Struct(structData)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, el)
		}
		return context.Map{"error": errors}
	}
	return nil
}

func CreatePushNotificationHandler(ctx iris.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print(utils.ReadXml("01"))
		return
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(context.Map{"response": nil, "error": utils.ReadXml("02")})
		log.Print(err)
		return
	}

	ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
	var pushNotification []CreatePushNotification
	err = json.Unmarshal(body, &pushNotification)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(context.Map{"response": nil, "error": utils.ReadXml("03")})
		log.Print(err)
		return
	}

	var pushMessageResponse []PushMessageResponse

	for i := 0; i < len(pushNotification); i++ {
		data := pushNotification[i]

		validation := ValidateStruct(data)
		if validation != nil {
			val := validation["error"].([]IError)
			response := PushMessageResponse{
				Status: iris.StatusBadRequest,
				Error:  fmt.Sprintf("%+v", val),
			}
			pushMessageResponse = append(pushMessageResponse, response)
			log.Print(validation)
			continue
		}

		switch data.Platform {
		case "iOS":
			authKey, err := token.AuthKeyFromFile(os.Getenv("P8_LOCATION"))
			if err != nil {
				log.Println("token error: ", err)
				response := PushMessageResponse{
					Status: iris.StatusOK,
					Error:  utils.ReadXml("04"),
				}
				pushMessageResponse = append(pushMessageResponse, response)
				continue
			}

			token2 := &token.Token{
				AuthKey: authKey,
				// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
				KeyID: data.KeyId,
				// TeamID from developer account (View Account -> Membership)
				TeamID: data.TeamId,
			}

			payload := payload2.NewPayload().AlertTitle(data.Title).Alert(data.Message)
			fmt.Println("payload")
			fmt.Println(payload)

			notification := &apns2.Notification{
				DeviceToken: data.DeviceToken,
				Topic:       os.Getenv("IOS_PACKAGE_NAME"),
				Payload:     payload,
			}

			client := apns2.NewTokenClient(token2).Production()
			res, err := client.Push(notification)
			if err != nil {
				log.Print("Send Push Message Error: ")
				log.Print(err)
				response := PushMessageResponse{
					Status: iris.StatusOK,
					Error:  err.Error(),
				}
				pushMessageResponse = append(pushMessageResponse, response)
				continue
			}

			response := PushMessageResponse{
				Status: res.StatusCode,
				Id:     res.ApnsID,
			}
			pushMessageResponse = append(pushMessageResponse, response)
			fmt.Println("pushMessageResponse")
			fmt.Println(pushMessageResponse)

			savePushNotification(data, res.ApnsID)
			break

		case "Android":
			opt := option.WithCredentialsFile(os.Getenv("ANDROID_JSON_LOCATION"))

			app, err := firebase.NewApp(contx.Background(), nil, opt)
			if err != nil {
				log.Println("new firebase app Error: ")
				response := PushMessageResponse{
					Status: iris.StatusOK,
					Error:  err.Error(),
				}
				pushMessageResponse = append(pushMessageResponse, response)
				continue
			}

			fcmClient, err := app.Messaging(contx.Background())
			if err != nil {
				log.Println("messaging Error: ")
				ctx.StatusCode(iris.StatusOK)
				response := PushMessageResponse{
					Status: iris.StatusOK,
					Error:  err.Error(),
				}
				pushMessageResponse = append(pushMessageResponse, response)
				continue
			}

			fmt.Println("data")
			fmt.Println(data)
			res, err := fcmClient.Send(contx.Background(), &messaging.Message{
				Notification: &messaging.Notification{
					Title: data.Title,
					Body:  data.Message,
				},
				Token: data.DeviceToken, // a token that you received from a client
			})

			if err != nil {
				log.Println("RESPONSE ERROR: ", err)
				response := PushMessageResponse{
					Status: iris.StatusOK,
					Error:  err.Error(),
				}
				pushMessageResponse = append(pushMessageResponse, response)
				continue
			}

			fmt.Println("response")
			fmt.Println(res)
			responseSplit := strings.Split(res, "/")

			response := PushMessageResponse{
				Status: 200,
				Id:     responseSplit[3],
			}
			pushMessageResponse = append(pushMessageResponse, response)

			fmt.Println("pushMessageResponse")
			fmt.Println(pushMessageResponse)

			savePushNotification(data, responseSplit[3])
			break
		}
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(context.Map{"response": pushMessageResponse, "error": nil})
	return
}

func savePushNotification(pushNotification CreatePushNotification, sendId string) {
	theTime := time.Now()
	dateData := theTime.Format("2006-01-02")
	timeData := theTime.Format("15:04")

	newPushNotification := model.PushNotification{
		Title:             pushNotification.Title,
		Message:           pushNotification.Message,
		DeviceTokenId:     pushNotification.DeviceTokenId,
		SendId:            sendId,
		OperativeSystemId: repository.GetIdOperativeSystem(pushNotification.Platform),
		OperativeSystem:   model.OperativeSystems{},
		Date:              dateData,
		Time:              timeData,
	}

	code, insertPushNotification, err := repository.InsertPushNotification(newPushNotification)
	if err != nil {
		log.Print(err)
		log.Print(code)
		return
	}
	log.Print("Push Message insertado de manera exitosa: ")
	log.Print(insertPushNotification)
}
