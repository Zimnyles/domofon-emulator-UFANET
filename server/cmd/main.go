package main

import (
	"context"
	"domofonEmulator/config"
	"domofonEmulator/pkg/database"
	"domofonEmulator/pkg/logger"
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/session"
	"domofonEmulator/server/internal/api"
	auth "domofonEmulator/server/internal/authHandler"
	intercom "domofonEmulator/server/internal/intercomHandler"
	mainhandler "domofonEmulator/server/internal/mainHandler"
	mqttserver "domofonEmulator/server/mqttServer"

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
	serverApp.Static("/server/web/public", "./server/web/public")
	serverApp.Static("/server/web/static", "./server/web/static")

	//mqtt
	mqttConfig := config.NewMQTTConfig()
	mqttServerRepository := mqttserver.NewMqttRepository(databasePool, logger)
	mqttServer, err := mqttserver.Connect(*mqttConfig, logger, *mqttServerRepository)
	if err != nil {
		logger.Fatal().Err(err).Msg("Client cannot connect to mqqtt")
	}
	defer mqttServer.Disconnect()

	//middlewares
	serverApp.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))
	serverApp.Use(recover.New())
	serverApp.Use(middleware.AuthMiddleware(store))

	//Repositories
	authRepository := auth.NewAuthRepository(databasePool, logger)
	intercomRepository := intercom.NewIntercomRepository(databasePool, logger)
	mainRepository := mainhandler.NewMainRepository(databasePool, logger)
	apiRepository := api.NewApiRepository(databasePool, logger)

	//Services
	authService := auth.NewAuthService(logger, *mqttServer, authRepository)
	intercomService := intercom.NewIntercomService(logger, *mqttServer, intercomRepository)
	mainService := mainhandler.NewMainService(logger, *mqttServer, intercomRepository)
	apiService := api.NewApiService(logger, *mqttServer, apiRepository)

	//Hadlers
	auth.NewHandler(serverApp, logger, *mqttServer, authService, authRepository, store)
	intercom.NewHandler(serverApp, logger, *mqttServer, intercomService, intercomRepository, store)
	mainhandler.NewHandler(serverApp, logger, *mqttServer, mainService, mainRepository, store)
	api.NewHandler(serverApp, logger, *mqttServer, apiService)

	go mqttServer.ListenForIntercomCreations(context.Background())
	go mqttServer.ListenForIntercomConnections(context.Background())
	go mqttServer.ListenForIntercomPowerOnOff(context.Background())
	go mqttServer.ListenForDoorControl(context.Background())
	go mqttServer.ListenForCalls(context.Background())
	go mqttServer.MonitorIntercomStatus(context.Background())

	serverApp.Listen(":3031")
}
