package home

import (
	"domofonEmulator/client/models"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type HomeService struct {
	logger     *zerolog.Logger
	mqqtClient mqtt.Client
}

func NewHomeService(logger *zerolog.Logger, mqqtClient mqtt.Client) *HomeService{
	return &HomeService{
		logger: logger,
		mqqtClient: mqqtClient,
	}
}

func(s *HomeService) NewIntercom(mac string, adress string, numofapartments int, topic string) error {
	newIntercom := models.NewIntercomProperties{
		MAC:                mac,
		Address:             adress,
		NumberOfApartments: numofapartments,
	}

	payload, err := json.Marshal(newIntercom)
	if err != nil {
		return fmt.Errorf("failed to marshal domofon data: %w", err)
	}

	token := s.mqqtClient.Publish(topic, 0, false, payload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("publish error: %w", token.Error())
	}

	return nil
}
