package intercom

import (
	"domofonEmulator/client/models"
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/storage"
	"domofonEmulator/client/web/views/components"
	"domofonEmulator/client/web/views/pages"
	"domofonEmulator/pkg/tadapter"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type IntercomHandler struct {
	router         fiber.Router
	logger         *zerolog.Logger
	mqqtClient     mqttclient.Client
	service        IIntercomService
	sessionStorage *storage.SessionStorage
	wsClients      map[string]map[*websocket.Conn]bool
}

type IIntercomService interface{}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqttclient.Client, service IIntercomService, sessionStorage *storage.SessionStorage) {
	h := &IntercomHandler{
		router:         router,
		logger:         logger,
		mqqtClient:     mqqtClient,
		service:        service,
		sessionStorage: sessionStorage,
		wsClients:      make(map[string]map[*websocket.Conn]bool),
	}

	h.router.Get("/intercom/:id", h.connectIntercome)
	h.router.Get("/ws/intercom/:id", websocket.New(h.websocketHandler))

	openTopic := "intercom/opendoor/+"
	err := mqqtClient.Subscribe(openTopic, func(payload []byte) {
		h.handleMqttMessage(payload)
	})
	if err != nil {
		logger.Error().Err(err).Msg("cannot subscribe to mqtt topic")
	}

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

func (h *IntercomHandler) handleMqttMessage(payload []byte) {
	var intercom models.Intercom
	if err := json.Unmarshal(payload, &intercom); err != nil {
		h.logger.Error().Err(err).Msg("mqtt message parsing error")
		return
	}

	conns, ok := h.wsClients[strconv.Itoa(intercom.ID)]
	if !ok {
		return
	}

	dataToSend, err := json.Marshal(intercom)
	if err != nil {
		h.logger.Error().Err(err).Msg("marshal data error")
		return
	}

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, dataToSend); err != nil {
			h.logger.Error().Err(err).Msg("ws sending error")
			conn.Close()
			delete(conns, conn)
		}
	}
}

func (h *IntercomHandler) websocketHandler(c *websocket.Conn) {
	intercomID := c.Params("id")

	if h.wsClients[intercomID] == nil {
		h.wsClients[intercomID] = make(map[*websocket.Conn]bool)
	}
	h.wsClients[intercomID][c] = true
	h.logger.Info().Msgf("ws connections success, id: %s", intercomID)

	defer func() {
		delete(h.wsClients[intercomID], c)
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
