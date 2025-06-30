package main

import (
	"domofonEmulator/config"
	"domofonEmulator/pkg/database"
	"domofonEmulator/pkg/logger"
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/session"

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

	//middlewares
	serverApp.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))
	serverApp.Use(recover.New())
	serverApp.Use(middleware.AuthMiddleware(store))

	serverApp.Listen(":3030")
}
