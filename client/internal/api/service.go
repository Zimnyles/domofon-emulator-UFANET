package api

import (
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

func (s *APIService) CreateNewIntercomRequest(mac string, address string, apartments int) (bool, string) {
	responseChan := make(chan struct {
		success bool
		id      int
		isNew   bool
		mac     string
		message string
	}, 1)

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
			case responseChan <- struct {
				success bool
				id      int
				isNew   bool
				mac     string
				message string
			}{
				success: false,
				message: "Ошибка обработки ответа сервера",
			}:
			default:
			}
			return
		}
		select {
		case responseChan <- struct {
			success bool
			id      int
			isNew   bool
			mac     string
			message string
		}{
			success: response.Success,
			id:      response.ID,
			isNew:   response.IsNew,
			mac:     response.Mac,
			message: response.Message,
		}:
		default:
		}
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to subscribe to MQTT topic")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору."
	}

	request := map[string]interface{}{
		"mac":        mac,
		"address":    address,
		"apartments": apartments,
	}

	payload, _ := json.Marshal(request)
	if err := s.mqqtClient.Publish("intercom/fromclient/create", payload); err != nil {
		s.logger.Error().Err(err).Msg("Failed to send creation request")
		return false, "Ошибка отправки запроса. Обратитесь к системному администратору."
	}

	s.logger.Info().
		Str("mac", mac).
		Str("address", address).
		Int("apartments", apartments).
		Msg("Sent intercom creation request")

	select {
	case response := <-responseChan:
		if response.isNew && response.success {
			return true, "Домофон успешно создан (ID: " + strconv.Itoa(response.id) + "). Используйте ID для подключния к устройству на главной странице."
		}
		if !response.isNew && response.success {
			return true, "Домофон с таким MAC адресом уже существует (ID: " + strconv.Itoa(response.id) + "). Используйте ID для подключния к устройству на главной странице."
		} else {
			return false, "Произошла неизвестная ошибка. Обратитесь к системному администратору." + response.message
		}
	case <-time.After(10 * time.Second):
		return false, "Превышено время ожидания ответа из-за ошибки на сервере. Обратитесь к системному администратору."
	}
}
