package home

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type HomeService struct {
	logger     *zerolog.Logger
	mqqtClient mqtt.Client
}

func NewHomeService(logger *zerolog.Logger, mqqtClient mqtt.Client) *HomeService {
	return &HomeService{
		logger:     logger,
		mqqtClient: mqqtClient,
	}
}
