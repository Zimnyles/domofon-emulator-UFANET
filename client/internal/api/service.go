package api

import (
	mqttclient "domofonEmulator/client/internal/mqttClient"
	"domofonEmulator/client/models"
	"domofonEmulator/config"
	"encoding/json"
	"fmt"

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

func (s *APIService) NewIntercom(mac string, adress string, numofapartments int, topic string) error {
	newIntercom := models.NewIntercomProperties{
		MAC:                mac,
		Address:            adress,
		NumberOfApartments: numofapartments,
	}

	payload, err := json.Marshal(newIntercom)
	if err != nil {
		return fmt.Errorf("failed to marshal domofon data: %w", err)
	}

	err = s.mqqtClient.Publish(topic, payload)

	if err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	return nil
}

// func (s *APIService) Run(ctx context.Context) {
// 	ticker := time.NewTicker(s.mqttConfig.StatusSendInterval)
// 	defer ticker.Stop()
// 	if err := s.publishStatuses(ctx); err != nil {
// 		s.logger.Printf("Initial publish failed: %v", err)
// 	}

// 	for {
// 		select {
// 		case <-ticker.C:
// 			if err := s.publishStatuses(ctx); err != nil {
// 				s.logger.Printf("Publish failed: %v", err)
// 			}
// 		case <-ctx.Done():
// 			s.logger.Println("Stopping doorphone service")
// 			return
// 		}
// 	}
// }

// func (s *APIService) publishStatuses(ctx context.Context) error {
// 	doorphones, err := s.repository.GetAllActiveIntercoms(ctx)
// 	if err != nil {
// 		s.logger.Info().Msg("Cannot get intercoms")
// 	}

// 	for _, d := range doorphones {
// 		topic := fmt.Sprintf("doorphones/%s/status", d.MAC)
// 		if err := s.mqqtClient.Publish(topic, d, 1); err != nil {
// 			s.logger.Printf("Failed to publish for doorphone %d: %v", d.ID, err)
// 			continue
// 		}
// 	}

// 	return nil
// }
