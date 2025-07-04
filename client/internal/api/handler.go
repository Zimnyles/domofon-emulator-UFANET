package api

import (
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/web/views/components"
	"domofonEmulator/pkg/tadapter"
	"strconv"

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
	CreateNewIntercomRequest(mac string, address string, apartments int) (bool, string)
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqttclient.Client, apiService IApiService) {
	h := &ApiHandler{
		router:     router,
		logger:     logger,
		mqqtClient: mqqtClient,
		apiService: apiService,
	}
	h.router.Post("/api/createIntercom", h.apiCreateIntercom)

}


func (h *ApiHandler) apiCreateIntercom(c *fiber.Ctx) error {
	mac := c.FormValue("mac")
	address := c.FormValue("address")
	apartments, _ := strconv.Atoi(c.FormValue("apartments"))

	if len(mac) > 17 {
		component := components.NewIntercomResponse("Длина MAC адреса не может превышать 17 символов!")
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	if apartments > 3800 {
		component := components.NewIntercomResponse("Количество квартир не может превышать 3800!")
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	if len(address) > 100 {
		component := components.NewIntercomResponse("Адрес не может первышать 100 символов!")
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	isSuccess, message := h.apiService.CreateNewIntercomRequest(mac, address, apartments)
	if isSuccess {
		h.logger.Info().
			Msg("Intercom mac is registered on the server")
	} else {
		h.logger.Info().
			Msg("Intercom msc is not registered on the server")
	}

	component := components.NewIntercomResponse(message)
	return tadapter.Render(c, component, fiber.StatusOK)
}
