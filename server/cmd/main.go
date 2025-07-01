package main

import (
	"domofonEmulator/config"
	"domofonEmulator/pkg/database"
	"domofonEmulator/pkg/logger"
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/session"
	"domofonEmulator/server/mqtt"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	//configs init
	config.Init()

	//logger
	loggerConfig := config.NewLogConfig()
	logger := logger.NewLogger(loggerConfig)

	//database
	databaseConfig := config.NewDBConfig()
	databasePool := database.CreateDataBasePool(databaseConfig, logger)
	defer databasePool.Close()

	//sessions
	sessionStorage := session.NewSessionStorage(databasePool)
	store := session.NewSession(sessionStorage)

	//app init
	serverApp := fiber.New()
	serverApp.Static("/web/static", "./web/static")

	//mqtt
	mqqtConfig := config.NewMQTTConfig()
	mqqtClient, err := mqtt.Connect(mqqtConfig.Broker)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MQTT broker")
		return
	}

	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		mqqtClient.Publish("test", 0, false, t.String())
	}

	//middlewares
	serverApp.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))
	serverApp.Use(recover.New())
	serverApp.Use(middleware.AuthMiddleware(store))

	serverApp.Listen(":3031")
}
