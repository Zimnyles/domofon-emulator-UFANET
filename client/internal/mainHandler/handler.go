package home

import (
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type HomeHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	mqqtClient mqttclient.Client
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqttclient.Client) {
	h := &HomeHandler{
		router:     router,
		logger:     logger,
		mqqtClient: mqqtClient,
	}

	h.router.Get("/", h.home)
	h.router.Get("/create", h.createIntercome)

}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	component := pages.HomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *HomeHandler) createIntercome(c *fiber.Ctx) error {
	component := pages.CreateIntercomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}
