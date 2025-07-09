package intercom

import (
	"domofonEmulator/server/models"
	mqttserver "domofonEmulator/server/mqttServer"

	"github.com/rs/zerolog"
)

type IntercomService struct {
	logger     *zerolog.Logger
	mqttServer mqttserver.Server
	repository IIntercomRepository
}

func NewIntercomService(logger *zerolog.Logger, mqttServer mqttserver.Server, repository IIntercomRepository) *IntercomService {
	return &IntercomService{
		logger:     logger,
		mqttServer: mqttServer,
		repository: repository,
	}
}

func (s *IntercomService) GetIntercomDataById(id int) (models.Intercom, error) {
	return s.repository.GetIntercomById(id)
}
