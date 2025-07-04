package intercom

import (
	"domofonEmulator/client/models"
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/storage"
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"
	"encoding/json"
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
	sess, err := h.sessionStorage.GetSession(c)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get session")
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	raw := sess.Get("intercom_data")
	jsonStr, ok := raw.(string)
	if !ok {
		h.logger.Warn().Msg("intercom_data not found or not a string")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid session data")
	}

	var intercomData models.Intercom
	if err := json.Unmarshal([]byte(jsonStr), &intercomData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal intercom data from session")
		return c.Status(fiber.StatusInternalServerError).SendString("Corrupted session data")
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
