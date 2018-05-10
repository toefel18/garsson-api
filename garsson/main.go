package main

import (
	"github.com/toefel18/garsson-api/garsson/logging"
	"github.com/toefel18/garsson-api/garsson/api"
	"os"
	"github.com/toefel18/garsson-api/garsson/db"
	"github.com/sirupsen/logrus"
)

//docker run --name loading-service-postgres -p 5434:5432 -e POSTGRES_USER=loadingservice -e POSTGRES_PASSWORD=loadingservice -d postgres
var ConnectionString = envOrDefault("CONNECTION_STRING", "postgres://garsson:garsson@localhost:5432/garsson?sslmode=disable")

func main() {
	logging.ConfigureDefault()
	logrus.Info("Starting Garsson")
	dao, err := db.NewDao(ConnectionString)
	if err != nil {
		logrus.WithError(err).Error("Invalid connection string")
		return
	}
	dao.WaitTillAvailable()
	api.Publish(dao)
}

func envOrDefault(key, defaultValue string) string {
	if val, present := os.LookupEnv(key); present {
		return val
	}
	return defaultValue
}