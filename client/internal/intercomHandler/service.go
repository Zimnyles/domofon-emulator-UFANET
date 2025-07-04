package intercom

import (
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/config"

	"github.com/rs/zerolog"
)

type IntercomService struct {
	logger     *zerolog.Logger
	mqqtClient mqttclient.Client
	mqttConfig config.MQTTConfig
}

type IIntercomRepository interface {
}

func NewIntercomService(logger *zerolog.Logger, mqqtClient mqttclient.Client, mqttConfig config.MQTTConfig) *IntercomService {
	return &IntercomService{
		logger:     logger,
		mqqtClient: mqqtClient,
		mqttConfig: mqttConfig,
	}
}
