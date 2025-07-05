package api

import (
	"domofonEmulator/client/models"
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/config"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type APIService struct {
	logger     *zerolog.Logger
	mqqtClient mqttclient.Client
	mqttConfig config.MQTTConfig
}

func NewAPIService(logger *zerolog.Logger, mqqtClient mqttclient.Client, mqttConfig config.MQTTConfig) *APIService {
	return &APIService{
		logger:     logger,
		mqqtClient: mqqtClient,
		mqttConfig: mqttConfig,
	}
}

func (s *APIService) PowerIntercomRequest(id int, action string) (bool, string, *models.Intercom) {
	responseChan := make(chan models.IntercomPowerOnOffResponse, 1)
	responseTopic := fmt.Sprintf("intercom/fromserver/power/%d", id)

	err := s.mqqtClient.Subscribe(responseTopic, func(payload []byte) {
		var response models.IntercomConnectResponse
		if err := json.Unmarshal(payload, &response); err != nil {
			s.logger.Error().Err(err).Str("payload", string(payload)).Msg("Failed to parse response")
			select {
			case responseChan <- models.IntercomPowerOnOffResponse{
				Success: false,
				Message: "Ошибка обработки ответа сервера",
			}:
			default:
			}
			return
		}
		select {
		case responseChan <- models.IntercomPowerOnOffResponse{
			Success:  response.Success,
			Message:  response.Message,
			Intercom: response.Intercom,
		}:
		default:
		}
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to subscribe to MQTT topic")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору.", nil
	}
	request := map[string]interface{}{
		"action": action,
		"id":     id,
	}

	payload, _ := json.Marshal(request)
	if err := s.mqqtClient.Publish("intercom/fromclient/power", payload); err != nil {
		s.logger.Error().Err(err).Msg("Failed to send power on/off request")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору.", nil
	}

	select {
	case response := <-responseChan:
		intercomData := models.Intercom{
			ID:                 response.ID,
			MAC:                response.MAC,
			IntercomStatus:     response.IntercomStatus,
			DoorStatus:         response.DoorStatus,
			Address:            response.Address,
			NumberOfApartments: response.NumberOfApartments,
			IsCalling:          response.IsCalling,
			CreatedAt:          response.CreatedAt,
			UpdatedAt:          response.UpdatedAt,
		}
		if response.Success {
			return true, "Запрос успешно доставлен и обработан сервером", &intercomData
		} else {
			return false, "Произошла неизвестная ошибка. Обратитесь к системному администратору" + response.Message, &intercomData
		}
	case <-time.After(10 * time.Second):
		return false, "Превышено время ожидания ответа из-за ошибки на сервере. Обратитесь к системному администратору", nil
	}
}

func (s *APIService) CreateNewIntercomConnectionRequest(id int) (bool, string, *models.Intercom) {
	responseChan := make(chan models.IntercomConnectResponse, 1)
	responseTopic := fmt.Sprintf("intercom/fromserver/connect/%d", id)

	err := s.mqqtClient.Subscribe(responseTopic, func(payload []byte) {
		var response models.IntercomConnectResponse
		if err := json.Unmarshal(payload, &response); err != nil {
			s.logger.Error().Err(err).Str("payload", string(payload)).Msg("Failed to parse response")
			select {
			case responseChan <- models.IntercomConnectResponse{
				Success: false,
				Message: "Ошибка обработки ответа сервера",
			}:
			default:
			}
			return
		}
		select {
		case responseChan <- models.IntercomConnectResponse{
			Success:  response.Success,
			Message:  response.Message,
			Intercom: response.Intercom,
		}:
		default:
		}
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to subscribe to MQTT topic")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору", nil
	}
	request := map[string]interface{}{
		"id": id,
	}

	payload, _ := json.Marshal(request)
	if err := s.mqqtClient.Publish("intercom/fromclient/connect", payload); err != nil {
		s.logger.Error().Err(err).Msg("Failed to send creation request")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору", nil
	}

	select {
	case response := <-responseChan:
		intercomData := models.Intercom{
			ID:                 response.ID,
			MAC:                response.MAC,
			IntercomStatus:     response.IntercomStatus,
			DoorStatus:         response.DoorStatus,
			Address:            response.Address,
			NumberOfApartments: response.NumberOfApartments,
			IsCalling:          response.IsCalling,
			CreatedAt:          response.CreatedAt,
			UpdatedAt:          response.UpdatedAt,
		}
		if response.Success {
			return true, "", &intercomData
		}
		if !response.Success && response.Message == "cannot find intercom by id" {
			return true, "Не удалось найти домофон по id", &intercomData
		} else {
			return false, "Произошла неизвестная ошибка. Обратитесь к системному администратору" + response.Message, &intercomData
		}
	case <-time.After(10 * time.Second):
		return false, "Превышено время ожидания ответа из-за ошибки на сервере. Обратитесь к системному администратору", nil
	}
}

func (s *APIService) CreateNewIntercomRequest(mac string, address string, apartments int) (bool, string) {
	responseChan := make(chan models.CreateIntercomResponse, 1)

	responseTopic := fmt.Sprintf("intercom/fromserver/%s", mac)

	err := s.mqqtClient.Subscribe(responseTopic, func(payload []byte) {
		var response struct {
			Success bool   `json:"success"`
			ID      int    `json:"id"`
			IsNew   bool   `json:"is_new"`
			Mac     string `json:"mac"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(payload, &response); err != nil {
			s.logger.Error().Err(err).Str("payload", string(payload)).Msg("Failed to parse response")
			select {
			case responseChan <- models.CreateIntercomResponse{
				Success: false,
				Message: "Ошибка обработки ответа сервера",
			}:
			default:
			}
			return
		}
		select {
		case responseChan <- models.CreateIntercomResponse{
			Success: response.Success,
			ID:      response.ID,
			IsNew:   response.IsNew,
			Mac:     response.Mac,
			Message: response.Message,
		}:
		default:
		}
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to subscribe to MQTT topic")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору"
	}

	request := map[string]interface{}{
		"mac":        mac,
		"address":    address,
		"apartments": apartments,
	}

	payload, _ := json.Marshal(request)
	if err := s.mqqtClient.Publish("intercom/fromclient/create", payload); err != nil {
		s.logger.Error().Err(err).Msg("Failed to send creation request")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору"
	}

	s.logger.Info().
		Str("mac", mac).
		Str("address", address).
		Int("apartments", apartments).
		Msg("Sent intercom creation request")

	select {
	case response := <-responseChan:
		if response.IsNew && response.Success {
			return true, "Домофон успешно создан (ID: " + strconv.Itoa(response.ID) + "). Используйте ID для подключния к устройству на главной странице"
		}
		if !response.IsNew && response.Success {
			return true, "Домофон с таким MAC адресом уже существует (ID: " + strconv.Itoa(response.ID) + "). Используйте ID для подключния к устройству на главной странице"
		} else {
			return false, "Произошла неизвестная ошибка. Обратитесь к системному администратору" + response.Message
		}
	case <-time.After(10 * time.Second):
		return false, "Превышено время ожидания ответа из-за ошибки на сервере. Обратитесь к системному администратору"
	}
}
