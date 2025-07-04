package api

import (
	"domofonEmulator/client/models"
	"domofonEmulator/client/storage"
	"domofonEmulator/client/web/views/components"
	"domofonEmulator/pkg/tadapter"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type ApiHandler struct {
	router         fiber.Router
	logger         *zerolog.Logger
	apiService     IApiService
	sessionStorage *storage.SessionStorage
}

type IApiService interface {
	CreateNewIntercomRequest(mac string, address string, apartments int) (bool, string)
	CreateNewIntercomConnectionRequest(id int) (bool, string, *models.Intercom)
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, apiService IApiService, sessionStorage *storage.SessionStorage) {
	h := &ApiHandler{
		router:         router,
		logger:         logger,
		apiService:     apiService,
		sessionStorage: sessionStorage,
	}
	h.router.Post("/api/createIntercom", h.apiCreateIntercom)
	h.router.Post("/api/connect", h.apiConnectToIntercom)

}

func (h *ApiHandler) apiConnectToIntercom(c *fiber.Ctx) error {
	intercomID, _ := strconv.Atoi(c.FormValue("intercomID"))
	if intercomID > 99999 {
		component := components.ConnectIntercomResponse("ID не может превышать 99999!")
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	h.logger.Info().Int("new connection request to intercom with id:", intercomID)

	isSuccess, message, intercomData := h.apiService.CreateNewIntercomConnectionRequest(intercomID)

	if message == "Не удалось найти домофон по id" {
		component := components.ConnectIntercomResponse(message)
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	if isSuccess && message == "" {
		sess, err := h.sessionStorage.GetSession(c)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to get session")
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}
		jsonData, err := json.Marshal(intercomData)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to marshal intercom data")
			return c.Status(fiber.StatusInternalServerError).SendString("Session serialization error")
		}

		sess.Set("intercom_data", string(jsonData))

		if err := h.sessionStorage.SaveSession(sess); err != nil {
			h.logger.Error().Err(err).Msg("Failed to save session")
			return c.Status(fiber.StatusInternalServerError).SendString("Session save error")
		}
		c.Set("HX-Redirect", fmt.Sprintf("/intercom/%d", intercomData.ID))
		return c.SendStatus(fiber.StatusNoContent)
	}

	component := components.ConnectIntercomResponse(message)
	return tadapter.Render(c, component, fiber.StatusOK)
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
