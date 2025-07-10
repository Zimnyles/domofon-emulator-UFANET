package mainhandler

import (
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/tadapter"
	mqttserver "domofonEmulator/server/mqttServer"
	"domofonEmulator/server/web/views/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type mainHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	service    IMainService
	repository IMainRepository
	store      *session.Store
}

type IMainRepository interface {
}

type IMainService interface {
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqttServer mqttserver.Server, service IMainRepository, repository IMainRepository, store *session.Store) {
	h := &mainHandler{
		router:     router,
		logger:     logger,
		service:    service,
		repository: repository,
		store:      store,
	}

	h.router.Get("/connect", middleware.AuthRequired(store), h.connectToIntercom)
	h.router.Get("/create", middleware.AuthRequired(store), h.createIntercome)

}

func (h *mainHandler) connectToIntercom(c *fiber.Ctx) error {
	component := pages.HomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *mainHandler) createIntercome(c *fiber.Ctx) error {
	component := pages.CreateIntercomePage()
	return tadapter.Render(c, component, fiber.StatusOK)
}
