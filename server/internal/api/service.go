package api

import (
	"domofonEmulator/server/models"
	mqttserver "domofonEmulator/server/mqttServer"
	"fmt"

	"github.com/rs/zerolog"
)

type ApiService struct {
	logger     *zerolog.Logger
	mqttServer mqttserver.Server
	repository IApiRepository
}

func NewApiService(logger *zerolog.Logger, mqttServer mqttserver.Server, repository IApiRepository) *ApiService {
	return &ApiService{
		logger:     logger,
		mqttServer: mqttServer,
		repository: repository,
	}
}

func (s *ApiService) CreateNewIntercom(mac, address string, apartments int) (int, bool, error) {
	return s.repository.CreateNewIntercom(mac, address, apartments)
}

func (s *ApiService) SendActualData(intercomData models.Intercom) error {
	topic := fmt.Sprintf("intercom/opendoor/%d", intercomData.ID)

	if err := s.mqttServer.Publish(topic, intercomData); err != nil {
		return fmt.Errorf("ошибка публикации MQTT: %w", err)
	}

	return nil
}

func (s *ApiService) GetIntercomDataById(id int) (models.Intercom, error) {
	return s.repository.GetIntercomById(id)
}

func (s *ApiService) SetIntercomDoorOpened(id int) error {
	return s.repository.SetIntercomDoorOpened(id)
}

func (s *ApiService) SetIntercomDoorClose(id int) error {
	return s.repository.SetIntercomDoorClose(id)
}