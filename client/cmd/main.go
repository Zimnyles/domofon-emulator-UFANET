package main

import (
	"domofonEmulator/client/internal/home"
	"domofonEmulator/config"
	"domofonEmulator/pkg/logger"
	"domofonEmulator/pkg/mqtt"
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

	//app init
	clientApp := fiber.New()
	clientApp.Static("/client/web/public", "./client/web/public")
	clientApp.Static("/client/web/static", "./client/web/static")

	mqqtConfig := config.NewMQTTConfig()
	mqqtClient := mqtt.Connect(mqqtConfig.Broker, logger)
	

	//middlewares
	clientApp.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))
	clientApp.Use(recover.New())

	//Services
	homeService := home.NewHomeService(logger, mqqtClient)

	//Hadlers
	home.NewHandler(clientApp, logger, mqqtClient, homeService)

	clientApp.Listen(":3030")
}
