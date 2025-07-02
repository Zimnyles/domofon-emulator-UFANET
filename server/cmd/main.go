package main

import (
	"domofonEmulator/config"
	"domofonEmulator/pkg/database"
	"domofonEmulator/pkg/logger"
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/mqtt"
	"domofonEmulator/pkg/session"
	"domofonEmulator/server/internal/auth"
	"domofonEmulator/server/internal/home"

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
	mqqtConfig := config.NewMQTTConfig()
	mqqtClient := mqtt.Connect(mqqtConfig.Broker, logger)

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
	homeService := home.NewHomeService(logger, mqqtClient)
	authService := auth.NewAuthService(logger, mqqtClient, authRepository)

	//Hadlers
	home.NewHandler(serverApp, logger, mqqtClient, homeService, homeRepository, store)
	auth.NewHandler(serverApp, logger, mqqtClient, authService, authRepository, store)

	serverApp.Listen(":3031")
}
