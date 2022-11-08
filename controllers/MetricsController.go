package controllers

import (
	"crypto/tls"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"log"
	"net/http"
	"os"
	"push/utils"
)

func GetMetricsHandler(ctx iris.Context) {
	mfChan := make(chan *dto.MetricFamily, 1024)
	transport, err := makeTransport()
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		err := godotenv.Load(".env")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(context.Map{"response": nil, "error": utils.ReadXml("01")})
			return
		}

		NodeExporterEndpoint := os.Getenv("EXPORTER_ENDPOINT")
		err = prom2json.FetchMetricFamilies(NodeExporterEndpoint, mfChan, transport)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(context.Map{"response": nil, "error": err})
			return
		}
	}()
	var result []*prom2json.Family
	for mf := range mfChan {
		result = append(result, prom2json.NewFamily(mf))
	}

	jsonText, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(context.Map{"response": nil, "error": utils.ReadXml("05")})
		return
	}
	ctx.Header("Content-Type", "application/json")
	if _, err := ctx.Write(jsonText); err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(context.Map{"response": nil, "error": utils.ReadXml("06")})
		return
	}
}

func makeTransport() (*http.Transport, error) {
	var transport *http.Transport
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return transport, nil
}
