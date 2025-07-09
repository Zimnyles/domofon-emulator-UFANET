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
	PowerIntercomRequest(id int, action string) (bool, string, *models.Intercom)
	OpenDoorRequest(id int, apartment int) (bool, string, *models.Intercom)
	CloseDoorRequest(id int) (bool, string, *models.Intercom)
	CallRequest(id int, apartment int) (bool, string, *models.Intercom)
	EndCallRequest(id int) (bool, string, *models.Intercom)
	SendActualData(intercomData models.Intercom) error
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, apiService IApiService, sessionStorage *storage.SessionStorage) {
	h := &ApiHandler{
		router:         router,
		logger:         logger,
		apiService:     apiService,
		sessionStorage: sessionStorage,
	}
	h.router.Post("/api/callIntercom", h.apiCallIntercom)
	h.router.Post("/api/endcallIntercom", h.apiEndCallIntercom)

	h.router.Post("/api/openIntercom", h.apiOpenIntercom)
	h.router.Post("/api/closeIntercom", h.apiCloseIntercom)

	h.router.Post("/api/createIntercom", h.apiCreateIntercom)
	h.router.Post("/api/powerIntercom", h.apiPowerIntercom)

	h.router.Post("/api/connect", h.apiConnectToIntercom)

}

func (h *ApiHandler) apiCallIntercom(c *fiber.Ctx) error {
	apartment, err := strconv.Atoi(c.FormValue("call"))
	if err != nil {
		return h.renderError(c, "Неизвестная ошибка. Обратитесь к системному администратору")
	}
	activeIntercomData, err := h.sessionStorage.GetActiveIntercomData(c)
	if err != nil {
		h.logger.Info().Err(err).Msg("Cannot get active intercom data")
		return h.renderError(c, "Ошибка сервера. Обратитесь к системному администратору")
	}
	if activeIntercomData.IsCalling {
		return h.renderError(c, "Звонок уже идет. Завершите, чтобы начать новый")
	}
	if !activeIntercomData.IntercomStatus {
		return h.renderError(c, "Домофон выключен. Совершить звонок невозможно")
	}

	isSuccess, message, intercomData := h.apiService.CallRequest(activeIntercomData.ID, apartment)
	if !isSuccess {
		h.logger.Error().Err(err).Msg("Failed to get response from server")
		return h.renderError(c, message)
	}

	err = h.sessionStorage.SetActiveIntercomData(c, *intercomData)
	h.apiService.SendActualData(*intercomData)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update session data")
		return h.renderError(c, "Ошибка сохранения данных. Обратитесь к системному администратору")
	}
	fmt.Println(intercomData)

	return tadapter.Render(c, tadapter.RenderIntercomAndNotificationResponse(message, intercomData), fiber.StatusOK)
}

func (h *ApiHandler) apiEndCallIntercom(c *fiber.Ctx) error {
	activeIntercomData, err := h.sessionStorage.GetActiveIntercomData(c)
	if err != nil {
		h.logger.Info().Err(err).Msg("Cannot get active intercom data")
		return h.renderError(c, "Ошибка сервера. Обратитесь к системному администратору")
	}
	if !activeIntercomData.IsCalling {
		return h.renderError(c, "Звонок не идет")
	}
	if !activeIntercomData.IntercomStatus {
		return h.renderError(c, "Домофон выключен. Звонок не идет")
	}

	isSuccess, message, intercomData := h.apiService.EndCallRequest(activeIntercomData.ID)
	if !isSuccess {
		h.logger.Error().Err(err).Msg("Failed to get response from server")
		return h.renderError(c, message)
	}

	err = h.sessionStorage.SetActiveIntercomData(c, *intercomData)
	h.apiService.SendActualData(*intercomData)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update session data")
		return h.renderError(c, "Ошибка сохранения данных. Обратитесь к системному администратору")
	}

	return tadapter.Render(c, tadapter.RenderIntercomAndNotificationResponse(message, intercomData), fiber.StatusOK)
}

func (h *ApiHandler) apiCloseIntercom(c *fiber.Ctx) error {
	activeIntercomData, err := h.sessionStorage.GetActiveIntercomData(c)
	if err != nil {
		h.logger.Info().Err(err).Msg("Cannot get active intercom data")
		return h.renderError(c, "Ошибка сервера. Обратитесь к системному администратору")
	}

	if !activeIntercomData.DoorStatus {
		return h.renderError(c, "Дверь уже закрыта")
	}

	if !activeIntercomData.IntercomStatus {
		return h.renderError(c, "Домофон выключен, дверь не может быть закрыта")
	}

	isSuccess, message, intercomData := h.apiService.CloseDoorRequest(activeIntercomData.ID)

	if !isSuccess {
		h.logger.Error().Err(err).Msg("Failed to get response from server")
		return h.renderError(c, message)
	}

	err = h.sessionStorage.SetActiveIntercomData(c, *intercomData)
	h.apiService.SendActualData(*intercomData)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update session data")
		return h.renderError(c, "Ошибка сохранения данных. Обратитесь к системному администратору")
	}

	return tadapter.Render(c, tadapter.RenderIntercomAndNotificationResponse(message, intercomData), fiber.StatusOK)
}

func (h *ApiHandler) apiOpenIntercom(c *fiber.Ctx) error {
	apartment, err := strconv.Atoi(c.FormValue("opendoor"))
	if err != nil {
		return h.renderError(c, "Неизвестная ошибка. Обратитесь к системному администратору")
	}

	activeIntercomData, err := h.sessionStorage.GetActiveIntercomData(c)
	if err != nil {
		h.logger.Info().Err(err).Msg("Cannot get active intercom data")
		return h.renderError(c, "Ошибка сервера. Обратитесь к системному администратору")
	}

	if !activeIntercomData.IntercomStatus {
		return h.renderError(c, "Домофон выключен. Дверь уже открыта")
	}

	if activeIntercomData.DoorStatus {
		return h.renderError(c, "Дверь уже открыта")
	}

	isSuccess, message, intercomData := h.apiService.OpenDoorRequest(activeIntercomData.ID, apartment)

	if !isSuccess {
		h.logger.Error().Err(err).Msg("Failed to get response from server")
		return h.renderError(c, message)
	}

	err = h.sessionStorage.SetActiveIntercomData(c, *intercomData)
	h.apiService.SendActualData(*intercomData)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update session data")
		return h.renderError(c, "Ошибка сохранения данных. Обратитесь к системному администратору")
	}

	return tadapter.Render(c, tadapter.RenderIntercomAndNotificationResponse(message, intercomData), fiber.StatusOK)

}

func (h *ApiHandler) apiPowerIntercom(c *fiber.Ctx) error {
	action := c.FormValue("action")
	if action != "on" && action != "off" {
		return h.renderError(c, "Неизвестная ошибка. Обратитесь к системному администратору")
	}

	activeIntercomData, err := h.sessionStorage.GetActiveIntercomData(c)
	if err != nil {
		h.logger.Info().Err(err).Msg("Cannot get active intercom data")
	} else {
		h.logger.Info().Int("Intercom ID:", activeIntercomData.ID).Msg("New intercom power off/onn request from client")
	}

	isSuccess, message, intercomData := h.apiService.PowerIntercomRequest(activeIntercomData.ID, action)

	if !isSuccess {
		h.logger.Error().Err(err).Msg("Failed to update session data")
		return h.renderError(c, message)
	}

	err = h.sessionStorage.SetActiveIntercomData(c, *intercomData)
	h.apiService.SendActualData(*intercomData)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update session data")
		return h.renderError(c, "Ошибка сохранения данных. Обратитесь к системному администратору")
	}

	return tadapter.Render(c, tadapter.RenderIntercomAndNotificationResponse(message, intercomData), fiber.StatusOK)
}

func (h *ApiHandler) apiConnectToIntercom(c *fiber.Ctx) error {
	intercomID, _ := strconv.Atoi(c.FormValue("intercomID"))
	if intercomID > 99999 {
		component := components.ConnectIntercomResponse("ID не может превышать 99999!")
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	h.logger.Info().Int("New connection request to intercom with id:", intercomID)

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

func (h *ApiHandler) renderError(c *fiber.Ctx, message string) error {
	status := "error"
	component := components.IntercomControlResponse(message, status)
	return tadapter.Render(c, component, fiber.StatusOK)
}
