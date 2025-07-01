package home

import (
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type IntercomeHandler struct {
	router fiber.Router
	logger *zerolog.Logger
}

func NewHandler(router fiber.Router, logger *zerolog.Logger) {
	h := &IntercomeHandler{
		router: router,
		logger: logger,
	}

	h.router.Get("/:username", h.intercome)
 
}

func (h *IntercomeHandler) intercome(c *fiber.Ctx) error {

	component := pages.HomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}
