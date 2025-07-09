package api

import (
	"domofonEmulator/pkg/tadapter"
	"domofonEmulator/server/models"
	mqttserver "domofonEmulator/server/mqttServer"
	"domofonEmulator/server/web/views/components"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type ApiHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	mqqtServer mqttserver.Server
	apiService IApiService
}

type IApiService interface {
	CreateNewIntercom(mac, address string, apartments int) (int, bool, error)
	SendActualData(intercomData models.Intercom) error
	GetIntercomDataById(id int) (models.Intercom, error)
	SetIntercomDoorOpened(id int) error
	SetIntercomDoorClose(id int) error
}

type IApiRepository interface {
	CreateNewIntercom(mac, address string, apartments int) (int, bool, error)
	GetIntercomById(id int) (models.Intercom, error)
	SetIntercomDoorOpened(id int) error
	SetIntercomDoorClose(id int) error
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtServer mqttserver.Server, apiService IApiService) {
	h := &ApiHandler{
		router:     router,
		logger:     logger,
		mqqtServer: mqqtServer,
		apiService: apiService,
	}

	h.router.Post("/api/createIntercom", h.apiCreateIntercom)
	h.router.Post("/api/opendoorIntercom", h.apiOpenIntercomDoor)
	h.router.Post("/api/closedoorIntercom", h.apiCloseIntercomDoor)

}

func (h *ApiHandler) apiCloseIntercomDoor(c *fiber.Ctx) error {
	intercomID, _ := strconv.Atoi(c.FormValue("intercom_id"))
	err := h.apiService.SetIntercomDoorClose(intercomID)
	if err != nil {
		h.logger.Err(err)
	}
	intercomData, _ := h.apiService.GetIntercomDataById(intercomID)

	h.apiService.SendActualData(intercomData)
	return nil
}

func (h *ApiHandler) apiOpenIntercomDoor(c *fiber.Ctx) error {
	intercomID, _ := strconv.Atoi(c.FormValue("intercom_id"))
	err := h.apiService.SetIntercomDoorOpened(intercomID)
	if err != nil {
		h.logger.Err(err)
	}
	intercomData, _ := h.apiService.GetIntercomDataById(intercomID)

	h.apiService.SendActualData(intercomData)
	return nil
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

	intercomId, isNew, err := h.apiService.CreateNewIntercom(mac, address, apartments)

	if err != nil {
		h.logger.Info().
			Msg("Intercom mac is registered on the server")
		message := "Произошла неизвестная ошибка. Обратитесь к системному администратору"
		component := components.NewIntercomResponse(message)
		return tadapter.Render(c, component, fiber.StatusOK)
	} else {
		h.logger.Info().
			Msg("Intercom msc is not registered on the server")
	}

	if !isNew {
		message := "Домофон с таким MAC адресом уже существует (ID: " + strconv.Itoa(intercomId) + "). Используйте ID для подключния к устройству на главной странице"
		component := components.NewIntercomResponse(message)
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	message := "Домофон успешно создан (ID: " + strconv.Itoa(intercomId) + "). Используйте ID для подключния к устройству на главной странице"

	component := components.NewIntercomResponse(message)
	return tadapter.Render(c, component, fiber.StatusOK)
}
