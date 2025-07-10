package intercom

import (
	"domofonEmulator/pkg/middleware"
	"domofonEmulator/pkg/tadapter"
	"domofonEmulator/server/models"
	mqttserver "domofonEmulator/server/mqttServer"
	"domofonEmulator/server/web/views/pages"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type IntercomHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	service    IIntercomService
	mqttServer mqttserver.Server
	repository IIntercomRepository
	store      *session.Store
	wsClients  map[string]map[*websocket.Conn]bool
}

type IIntercomRepository interface {
	GetIntercomById(id int) (models.Intercom, error)
}

type IIntercomService interface {
	GetIntercomDataById(id int) (models.Intercom, error)
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqttServer mqttserver.Server, service IIntercomService, repository IIntercomRepository, store *session.Store) {
	h := &IntercomHandler{
		router:     router,
		logger:     logger,
		mqttServer: mqttServer,
		service:    service,
		repository: repository,
		store:      store,
		wsClients:  make(map[string]map[*websocket.Conn]bool),
	}

	h.router.Get("/intercom/:id", middleware.AuthRequired(store), h.intercomControl)
	h.router.Get("/ws/intercom/:id", middleware.AuthRequired(store), websocket.New(h.websocketHandler))
	h.router.Get("/ws/intercom/status/:id", middleware.AuthRequired(store), websocket.New(h.websocketHandler))
	h.router.Get("/ws/intercom/actualstatus/:id", middleware.AuthRequired(store), websocket.New(h.websocketHandler))
	h.router.Get("/ws/intercom/opendoor/:id", middleware.AuthRequired(store), websocket.New(h.websocketHandler))
	h.router.Get("/ws/intercom/activestatus/:id", middleware.AuthRequired(store), websocket.New(h.websocketHandler))

	h.router.Post("/intercom/redirect", middleware.AuthRequired(store), h.redirectToIntercom)

	topicStatus := "intercom/status/+"
	err := mqttServer.Subscribe(topicStatus, func(payload []byte) {
		h.handleMqttMessage("status", payload)
	})
	if err != nil {
		logger.Error().Err(err).Msg("cannot subscribe to mqtt topic")
	}

	topicActiveStatus := "intercom/activestatus/+"
	err = mqttServer.Subscribe(topicActiveStatus, func(payload []byte) {
		h.handleMqttMessage("activestatus", payload)
	})
	if err != nil {
		logger.Error().Err(err).Msg("cannot subscribe to mqtt topic")
	}

	topicActualStatus := "intercom/actualstatus/+"
	err = mqttServer.Subscribe(topicActualStatus, func(payload []byte) {
		h.handleMqttMessage("actualstatus", payload)
	})
	if err != nil {
		logger.Error().Err(err).Msg("cannot subscribe to mqtt topic")
	}

	topicActualStatusOpenDoor := "intercom/opendoor/+"
	err = mqttServer.Subscribe(topicActualStatusOpenDoor, func(payload []byte) {
		h.handleMqttMessage("actualstatus", payload)
	})
	if err != nil {
		logger.Error().Err(err).Msg("cannot subscribe to mqtt topic")
	}

}

func (h *IntercomHandler) handleMqttMessage(topicPrefix string, payload []byte) {
	var msg struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(payload, &msg); err != nil {
		h.logger.Error().Err(err).Msg("mqtt message parsing error")
		return
	}

	intercomID := strconv.Itoa(msg.ID)

	intercomData, err := h.service.GetIntercomDataById(msg.ID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get intercom data")
		return
	}

	wrappedMsg := map[string]interface{}{
		topicPrefix: intercomData,
	}
	wrappedData, err := json.Marshal(wrappedMsg)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal response")
		return
	}

	h.broadcastToClients(intercomID, wrappedData)
}

func (h *IntercomHandler) redirectToIntercom(c *fiber.Ctx) error {
	id := c.FormValue("intercomID")
	c.Response().Header.Add("Hx-Redirect", "/intercom/"+id)
	return c.Redirect("/", http.StatusOK)
}

func (h *IntercomHandler) intercomControl(c *fiber.Ctx) error {
	intercomID, _ := strconv.Atoi(c.Params("id"))

	intercomData, _ := h.service.GetIntercomDataById(intercomID)

	sess, err := h.store.Get(c)
	if err != nil {
		h.logger.Fatal().Err(err).Msg("Failed to get session store")
		panic(err)
	}

	userLogin := sess.Get("login").(string)

	component := pages.LiveIntercomPage(intercomData, userLogin)
	return tadapter.Render(c, component, fiber.StatusOK)

}

func (h *IntercomHandler) websocketHandler(c *websocket.Conn) {
	intercomID := c.Params("id")

	if h.wsClients[intercomID] == nil {
		h.wsClients[intercomID] = make(map[*websocket.Conn]bool)
	}
	h.wsClients[intercomID][c] = true
	h.logger.Info().Msgf("ws connection success, id: %s", intercomID)

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

func (h *IntercomHandler) broadcastToClients(intercomID string, message []byte) {
	conns, ok := h.wsClients[intercomID]
	if !ok {
		return
	}
	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			h.logger.Error().Err(err).Msg("ws sending error")
			conn.Close()
			delete(conns, conn)
		}
	}
}
