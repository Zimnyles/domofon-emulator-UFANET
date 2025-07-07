package intercom

import (
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/storage"
	"domofonEmulator/client/web/views/components"
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type IntercomHandler struct {
	router         fiber.Router
	logger         *zerolog.Logger
	mqqtClient     mqttclient.Client
	service        IIntercomService
	sessionStorage *storage.SessionStorage
}

type IIntercomService interface{}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqttclient.Client, service IIntercomService, sessionStorage *storage.SessionStorage) {
	h := &IntercomHandler{
		router:         router,
		logger:         logger,
		mqqtClient:     mqqtClient,
		service:        service,
		sessionStorage: sessionStorage,
	}

	h.router.Get("/intercom/:id", h.connectIntercome)

}

func (h *IntercomHandler) connectIntercome(c *fiber.Ctx) error {

	intercomData, err := h.sessionStorage.GetActiveIntercomData(c)
	if err != nil {
		component := components.ConnectIntercomResponse("Ошибка сервера. Обратитесь к системному администратору")
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	linkIDstring := c.Params("id")
	linkID, _ := strconv.Atoi(linkIDstring)

	if linkID != intercomData.ID {
		if c.Get("HX-Request") == "true" {
			c.Set("HX-Redirect", "/")
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Redirect("/", fiber.StatusSeeOther)
		}
	}

	component := pages.ControlInetcomPage(intercomData)
	return tadapter.Render(c, component, fiber.StatusOK)
}
