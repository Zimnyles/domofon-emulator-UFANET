package main

import (
	"domofonEmulator/client/internal/api"
	intercom "domofonEmulator/client/internal/intercomHandler"
	home "domofonEmulator/client/internal/mainHandler"
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/storage"
	"domofonEmulator/config"
	"domofonEmulator/pkg/database"
	"domofonEmulator/pkg/logger"

	// "time"

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

	//app init
	clientApp := fiber.New()
	clientApp.Static("/client/web/public", "./client/web/public")
	clientApp.Static("/client/web/static", "./client/web/static")

	//redis config
	redisConfig := config.NewRedisConfig()

	//session
	sessionStorage := storage.NewRedisStorage(*redisConfig)

	//mqtt
	mqttConfig := config.NewMQTTConfig()
	mqttClient, err := mqttclient.Connect(*mqttConfig, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Client cannot connect to mqqtt")
	}
	defer mqttClient.Disconnect()

	//middlewares
	clientApp.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))
	clientApp.Use(recover.New())

	//Repositories

	//Services
	intercomService := intercom.NewIntercomService(logger, *mqttClient, *mqttConfig)
	apiService := api.NewAPIService(logger, *mqttClient, *mqttConfig)

	//Hadlers
	home.NewHandler(clientApp, logger, *mqttClient)
	api.NewHandler(clientApp, logger, apiService, sessionStorage)
	intercom.NewHandler(clientApp, logger, *mqttClient, intercomService, sessionStorage)

	//Intercoms statuses sending
	// go intercomService.RunIntercomStatusSend()

	clientApp.Listen(":3030")
}
