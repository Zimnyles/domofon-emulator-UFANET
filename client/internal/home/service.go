package home

import (
	mqttclient "domofonEmulator/client/mqttClient"
	"domofonEmulator/client/models"
	"domofonEmulator/config"
	"sync"

	"github.com/rs/zerolog"
)

type IntercomService struct {
	logger          *zerolog.Logger
	mqqtClient      mqttclient.Client
	mqttConfig      config.MQTTConfig
	mu              sync.Mutex
	currentIntercom *models.Intercom
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

// func (s *IntercomService) SetupSubscriptions() {
// 	topic := "intercom/response"
// 	_ = s.mqqtClient.Subscribe(topic, func(payload []byte) {
// 		var resp models.RegistrationResponse
// 		if err := json.Unmarshal(payload, &resp); err != nil {
// 			s.logger.Error().Err(err).Msg("Failed to parse response")
// 			return
// 		}

// 		s.mu.Lock()
// 		defer s.mu.Unlock()

// 		if resp.Success && resp.Intercom != nil {
// 			s.currentIntercom = resp.Intercom
// 			s.logger.Info().Str("mac", resp.Intercom.MAC).Msg("Doorphone registered/connected")
// 		}
// 	})
// }

// func (s *IntercomService) RegisterOrConnect(doorphone *models.Intercom, action string) error {
//     req := models.RegistrationRequest{
//         Action:            action,
//         MAC:               doorphone.MAC,
//         Address:           doorphone.Address,
//         NumberOfApartments: doorphone.NumberOfApartments,
//     }

//     topic := "doorphones/register"
//     return s.mqqtClient.Publish(topic, req, 1)
// }

// func (s *IntercomService) GetCurrentDoorphone() *models.Intercom {
//     s.mu.Lock()
//     defer s.mu.Unlock()
//     return s.currentIntercom
// }

// func (s *IntercomService) RunIntercomStatusSend() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		s.logger.Info().Msg("Starting sending intercoms statuses")
// 		s.Run(ctx)
// 	}()

// 	sig := <-sigChan
// 	s.logger.Info().Any("Shutting down.Received signal:", sig)
// 	cancel()

// 	select {
// 	case <-time.After(2 * time.Second):
// 		s.logger.Info().Msg("Clean shutdown completed")
// 	case <-sigChan:
// 		s.logger.Info().Msg("Forced shutdown :(")
// 	}
// }
