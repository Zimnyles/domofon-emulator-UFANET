package mainhandler

import (
	mqttserver "domofonEmulator/server/mqttServer"

	"github.com/rs/zerolog"
)

type MainService struct {
	logger     *zerolog.Logger
	mqttServer mqttserver.Server
	repository IMainRepository
}

func NewMainService(logger *zerolog.Logger, mqttServer mqttserver.Server, repository IMainRepository) *MainService {
	return &MainService{
		logger:     logger,
		mqttServer: mqttServer,
		repository: repository,
	}
}
