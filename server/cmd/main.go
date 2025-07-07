package main

import (
	"context"
	"domofonEmulator/config"
	"domofonEmulator/pkg/database"
	"domofonEmulator/pkg/logger"
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/session"
	"domofonEmulator/server/internal/auth"
	"domofonEmulator/server/internal/home"
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
	homeRepository := home.NewHomeRepository(databasePool, logger)
	authRepository := auth.NewAuthRepository(databasePool, logger)

	//Services
	homeService := home.NewHomeService(logger, *mqttServer)
	authService := auth.NewAuthService(logger, *mqttServer, authRepository)

	//Hadlers
	home.NewHandler(serverApp, logger, *mqttServer, homeService, homeRepository, store)
	auth.NewHandler(serverApp, logger, *mqttServer, authService, authRepository, store)

	go mqttServer.ListenForIntercomCreations(context.Background())
	go mqttServer.ListenForIntercomConnections(context.Background())
	go mqttServer.ListenForIntercomPowerOnOff(context.Background())
	go mqttServer.ListenForDoorControl(context.Background())
	go mqttServer.ListenForCalls(context.Background())

	serverApp.Listen(":3031")
}
