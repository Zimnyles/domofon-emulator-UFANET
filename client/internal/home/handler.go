package home

import (
	mqttclient "domofonEmulator/client/internal/mqttClient"
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type HomeHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	mqqtClient mqttclient.Client
	service    IIntercomService
}

type IIntercomService interface {
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqttclient.Client, service IIntercomService) {
	h := &HomeHandler{
		router:     router,
		logger:     logger,
		mqqtClient: mqqtClient,
		service:    service,
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
