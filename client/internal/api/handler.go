package api

import (
	mqttclient "domofonEmulator/client/internal/mqttClient"
	"domofonEmulator/client/web/views/components"
	"domofonEmulator/pkg/tadapter"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type ApiHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	mqqtClient mqttclient.Client
	apiService IApiService
}

type IApiService interface {
	NewIntercom(mac string, adress string, numofapartments int, topic string) error
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqttclient.Client, apiService IApiService) {
	h := &ApiHandler{
		router:     router,
		logger:     logger,
		mqqtClient: mqqtClient,
		apiService: apiService,
	}

	// h.router.Post("/api/connect", h.connectIntercom)
	h.router.Post("/api/createIntercom", h.apiCreateIntercom)

}

type ConnectRequest struct {
	MAC        string `json:"mac"`
	Address    string `json:"address"`
	Apartments int    `json:"apartments"`
}

type IntercomResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (h *ApiHandler) apiCreateIntercom(c *fiber.Ctx) error {
	mac := c.FormValue("mac")
	address := c.FormValue("address")
	apartmentsString := c.FormValue("apartments")
	apartments, _ := strconv.Atoi(apartmentsString)

	responseChan := make(chan struct {
		success bool
		message string
		mac     string
	})
	defer close(responseChan)

	responseTopic := "intercom/create/response/" + mac
	err := h.mqqtClient.Subscribe(responseTopic, func(payload []byte) {
		var response struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
			Mac     string `json:"mac"`
		}
		if err := json.Unmarshal(payload, &response); err != nil {
			responseChan <- struct {
				success bool
				message string
				mac     string
			}{false, "Ошибка обработки ответа сервера",""}
			return
		}
		responseChan <- struct {
			success bool
			message string
			mac     string
		}{response.Success, response.Message, response.Mac}
	})

	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to subscribe to MQTT topic")
		component := components.NewIntercomResponse("Ошибка подключения к MQTT")
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	err = h.apiService.NewIntercom(mac, address, apartments, "intercom/create")
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
	select {
	case response := <-responseChan:
		if response.success {
			component := components.NewIntercomResponse("Домофон успешно создан: " + response.message)
			return tadapter.Render(c, component, fiber.StatusOK)
		} else {
			component := components.NewIntercomResponse("Ошибка: " + response.message)
			return tadapter.Render(c, component, fiber.StatusOK)
		}
	case <-time.After(30 * time.Second):
		component := components.NewIntercomResponse("Превышено время ожидания ответа от сервера")
		return tadapter.Render(c, component, fiber.StatusOK)
	}
}

func (h *ApiHandler) HandleIntercomAction(c *fiber.Ctx, action string) error {
	var req ConnectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(IntercomResponse{
			Success: false,
			Message: "Invalid request format",
		})
	}

	payload := map[string]interface{}{
		"action":     action,
		"mac":        req.MAC,
		"address":    req.Address,
		"apartments": req.Apartments,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal payload")
		return c.Status(fiber.StatusInternalServerError).JSON(IntercomResponse{
			Success: false,
			Message: "Internal server error",
		})
	}

	if err := h.mqqtClient.Publish("doorphones/"+action, jsonData); err != nil {
		h.logger.Error().Err(err).Msg("Failed to send MQTT message")
		return c.Status(fiber.StatusInternalServerError).JSON(IntercomResponse{
			Success: false,
			Message: "Failed to send request",
		})
	}

	return c.JSON(IntercomResponse{
		Success: true,
		Message: fmt.Sprintf("%s request sent", action),
	})
}
