package main

import (
	"domofonEmulator/client/internal/api"
	"domofonEmulator/client/internal/home"
	mqttclient "domofonEmulator/client/mqttClient"
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

	//mqtt
	mqqtConfig := config.NewMQTTConfig()
	mqqtClient, err := mqttclient.Connect(*mqqtConfig, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Client cannot connect to mqqtt")
	}
	defer mqqtClient.Disconnect()

	//middlewares
	clientApp.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))
	clientApp.Use(recover.New())

	//Repositories

	//Services
	intercomService := home.NewIntercomService(logger, *mqqtClient, *mqqtConfig)
	apiService := api.NewAPIService(logger, *mqqtClient, *mqqtConfig)

	//Hadlers
	home.NewHandler(clientApp, logger, *mqqtClient, intercomService)
	api.NewHandler(clientApp, logger, *mqqtClient, apiService)

	//Intercoms statuses sending
	// go intercomService.RunIntercomStatusSend()

	clientApp.Listen(":3030")
}
