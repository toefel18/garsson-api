package main

import (
    "encoding/hex"
    "fmt"
    "os"

    "github.com/toefel18/garsson-api/garsson/api"
    "github.com/toefel18/garsson-api/garsson/db"
    "github.com/toefel18/garsson-api/garsson/db/migration"
    "github.com/toefel18/garsson-api/garsson/log"
    "golang.org/x/crypto/sha3"
)

//docker run --name garsson-api-postgres -p 5432:5432 -e POSTGRES_USER=garsson -e POSTGRES_PASSWORD=garsson -d postgres
var ConnectionString = envOrDefault("CONNECTION_STRING", "postgres://garsson:garsson@localhost:5432/garsson?sslmode=disable")

func main() {
    doHash()
    log.ConfigureDefault()
    log.Info("Starting Garsson")
    dao, err := db.NewDao(ConnectionString)
    if err != nil {
        log.WithError(err).Error("Invalid connection string")
        return
    }
    dao.WaitTillAvailable()
    migration.MigrateDatabase(dao.NewSession())
    api := api.NewServer(dao)
    api.Start()
}

func envOrDefault(key, defaultValue string) string {
    if val, present := os.LookupEnv(key); present {
        return val
    }
    return defaultValue
}



