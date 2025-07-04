package home

import (
	mqttserver "domofonEmulator/server/mqttServer"

	"github.com/rs/zerolog"
)

type HomeService struct {
	logger     *zerolog.Logger
	mqttServer mqttserver.Server
}

func NewHomeService(logger *zerolog.Logger, mqttServer mqttserver.Server) *HomeService {
	return &HomeService{
		logger:     logger,
		mqttServer: mqttServer,
	}
}
