package home

import (
	"domofonEmulator/pkg/middleware"
	mqttserver "domofonEmulator/server/mqttServer"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type HomeHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	mqttServer mqttserver.Server
	service    IHomeService
	repository IHomeRepository
	store      *session.Store
}

type IHomeRepository interface {
}

type IHomeService interface {
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqttServer mqttserver.Server, service IHomeService, repository IHomeRepository, store *session.Store) {
	h := &HomeHandler{
		router:     router,
		logger:     logger,
		mqttServer: mqttServer,
		service:    service,
		repository: repository,
		store:      store,
	}

	h.router.Get("/", middleware.AuthRequired(h.store), h.home)

}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	return nil
}
