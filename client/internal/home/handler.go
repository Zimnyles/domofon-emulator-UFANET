package home

import (
	"domofonEmulator/client/web/views/components"
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type HomeHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	mqqtClient mqtt.Client
	service    IHomeService
}

type IHomeService interface {
	NewIntercom(mac string, adress string, numofapartments int, topic string) error
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqtt.Client, service IHomeService) {
	h := &HomeHandler{
		router:     router,
		logger:     logger,
		mqqtClient: mqqtClient,
		service:    service,
	}

	h.router.Get("/", h.home)
	h.router.Get("/create", h.createIntercome)

	h.router.Post("/api/createIntercom", h.apiCreateIntercom)

}

func (h *HomeHandler) home(c *fiber.Ctx) error {

	component := pages.HomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *HomeHandler) createIntercome(c *fiber.Ctx) error {

	component := pages.CreateIntercomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *HomeHandler) apiCreateIntercom(c *fiber.Ctx) error {
	mac := c.FormValue("mac")
	address := c.FormValue("address")
	apartmentsString := c.FormValue("apartments")
	apartments, err := strconv.Atoi(apartmentsString)
	if err != nil {
		h.logger.Fatal().Err(err).Msg("Failed to connect to MQTT broker")
	}
	err = h.service.NewIntercom(mac, address, apartments, "intercom/create")
	if err != nil {
		h.logger.Fatal().Err(err).Msg("Failed to create new intercom")
		component := components.NewIntercomResponse("Не удалось отправить запрос на сервер")
		return tadapter.Render(c, component, fiber.StatusOK)
	} else {
		h.logger.Info().
			Str("mac", mac).
			Str("address", address).
			Int("apartments", apartments).
			Str("mqtt_topic", "intercom/create").
			Msg("New domofon registered in system via MQTT")
	}
	component := components.NewIntercomResponse("Запрос отправлен на сервер")
	return tadapter.Render(c, component, fiber.StatusOK)
}
