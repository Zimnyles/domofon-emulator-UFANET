package mqttserver

import (
	"context"
	"domofonEmulator/config"
	"domofonEmulator/server/models"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type Server struct {
	client               mqtt.Client
	logger               *zerolog.Logger
	config               config.MQTTConfig
	mqqtServerRepository MqttServerRepository
}

func Connect(mqqtConfig config.MQTTConfig, logger *zerolog.Logger, mqqtServerRepository MqttServerRepository) (*Server, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqqtConfig.Broker)
	opts.SetClientID("server")
	opts.SetAutoReconnect(true)
	opts.SetResumeSubs(true)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MQTT broker")
		return nil, err

	} else {
		logger.Info().Msg("Client is connected to MQTT broker")
	}
	return &Server{
		client:               client,
		logger:               logger,
		mqqtServerRepository: mqqtServerRepository,
	}, nil
}

func (s *Server) Subscribe(topic string, handler func(payload []byte)) error {
	token := s.client.Subscribe(topic, byte(s.config.QOSLevel), func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Server) Disconnect() {
	s.client.Disconnect(250)
	s.logger.Info().Msg("MQTT client disconnected")
}

func (s *Server) ListenForIntercomPowerOnOff(ctx context.Context) {
	err := s.Subscribe("intercom/fromclient/power", func(payload []byte) {
		if s.mqqtServerRepository.Dbpool == nil {
			s.logger.Fatal().Msg("MQTT repository or DB pool is nil")
			return
		}
		var powerRequest struct {
			ID     int    `json:"id"`
			Action string `json:"action"`
		}
		if err := json.Unmarshal(payload, &powerRequest); err != nil {
			s.logger.Error().
				Err(err).
				Str("payload", string(payload)).
				Msg("Failed to parse intercom power on/off request message")
			return
		}
		s.logger.Info().
			Int("ID", powerRequest.ID).
			Str("Action", powerRequest.Action).
			Msg("Received new intercom power on/off change request")
		newStatus := powerRequest.Action == "on"
		err := s.mqqtServerRepository.UpdateIntercomStatus(powerRequest.ID, newStatus)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to update intecrom data")
			errorResponse := map[string]interface{}{
				"success": false,
				"message": "failed to update intercom status on server",
			}
			errTopic := fmt.Sprintf("intercom/fromserver/power/%d", powerRequest.ID)
			if err := s.Publish(errTopic, errorResponse); err != nil {
				s.logger.Error().Err(err).Msg("Failed to send error response")
			}
			return
		}
		updatedIntercom, err := s.mqqtServerRepository.GetIntercomByID(powerRequest.ID, s.logger)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to get intecrom data")
			errorResponse := map[string]interface{}{
				"success": false,
				"message": "failed to get intercom data on server",
			}
			errTopic := fmt.Sprintf("intercom/fromserver/power/%d", powerRequest.ID)
			if err := s.Publish(errTopic, errorResponse); err != nil {
				s.logger.Error().Err(err).Msg("Failed to send error response")
			}
			return
		}
		fullUpdateTopic := fmt.Sprintf("intercom/fromserver/power/%d", powerRequest.ID)
		fullResponse := map[string]interface{}{
			"success":              true,
			"message":              "ok",
			"id":                   updatedIntercom.ID,
			"mac_address":          updatedIntercom.MAC,
			"door_status":          updatedIntercom.DoorStatus,
			"intercom_status":      updatedIntercom.IntercomStatus,
			"address":              updatedIntercom.Address,
			"number_of_apartments": updatedIntercom.NumberOfApartments,
			"is_calling":           updatedIntercom.IsCalling,
			"created_at":           updatedIntercom.CreatedAt,
			"updated_at":           updatedIntercom.UpdatedAt,
		}
		if err := s.Publish(fullUpdateTopic, fullResponse); err != nil {
			s.logger.Error().
				Err(err).
				Str("topic", fullUpdateTopic).
				Msg("Failed to send intercom data")
			return
		}
		s.logger.Info().
			Int("ID", powerRequest.ID).
			Bool("NewStatus", newStatus).
			Msg("Successfully updated and broadcasted updated intercom data")
	})
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to subscribe to intercom power topic")
		return
	}
	<-ctx.Done()
	s.logger.Info().Msg("Stopping MQTT listener for intercom power state changes")
}

func (s *Server) ListenForIntercomConnections(ctx context.Context) {
	err := s.Subscribe("intercom/fromclient/connect", func(payload []byte) {
		if s.mqqtServerRepository.Dbpool == nil {
			s.logger.Fatal().Msg("MQTT repository or DB pool is nil!")
			return
		}

		var intercomConnectRequestData models.CreateIntercomConnectionRequset

		if err := json.Unmarshal(payload, &intercomConnectRequestData); err != nil {
			s.logger.Error().
				Err(err).
				Str("payload", string(payload)).
				Msg("Failed to parse intercom connection request message from client")
			return
		}
		s.logger.Info().
			Int("ID", intercomConnectRequestData.ID).
			Msg("Received new intercom connection request message request from client.")

		intercomData, err := s.mqqtServerRepository.GetIntercomByID(intercomConnectRequestData.ID, s.logger)

		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to get intecrom data")
			errorResponse := map[string]interface{}{
				"success": false,
				"message": "cannot find intercom by id",
			}
			errTopic := fmt.Sprintf("intercom/fromserver/connect/%d", intercomConnectRequestData.ID)
			if err := s.Publish(errTopic, errorResponse); err != nil {
				s.logger.Error().Err(err).Msg("Failed to send error response")
			}
			return
		}

		if intercomData.ID == 0 {
			s.logger.Error().Err(err).Msg("Failed to get intecrom data")
			errorResponse := map[string]interface{}{
				"success": false,
				"message": "cannot find intercom by id",
			}
			errTopic := fmt.Sprintf("intercom/fromserver/connect/%d", intercomConnectRequestData.ID)
			if err := s.Publish(errTopic, errorResponse); err != nil {
				s.logger.Error().Err(err).Msg("Failed to send error response")
			}
			return
		}

		responseTopic := fmt.Sprintf("intercom/fromserver/connect/%d", intercomData.ID)
		response := map[string]interface{}{
			"success":              true,
			"message":              "ok",
			"id":                   intercomData.ID,
			"mac_address":          intercomData.MAC,
			"door_status":          intercomData.DoorStatus,
			"address":              intercomData.Address,
			"number_of_apartments": intercomData.NumberOfApartments,
			"is_calling":           intercomData.IsCalling,
			"created_at":           intercomData.CreatedAt,
			"updated_at":           intercomData.UpdatedAt,
		}
		if err := s.Publish(responseTopic, response); err != nil {
			s.logger.Error().Err(err).Str("topic", responseTopic).Msg("Failed to send intercom connection response")
			return
		}
		s.logger.Info().Msg("Successfully sent intercom creation response")

	})
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to subscribe to intercom creation topic")
		return
	}

	<-ctx.Done()
	s.logger.Info().Msg("Stopping MQTT listener for intercom creations")

}

func (s *Server) ListenForIntercomCreations(ctx context.Context) {
	err := s.Subscribe("intercom/fromclient/create", func(payload []byte) {
		if s.mqqtServerRepository.Dbpool == nil {
			s.logger.Fatal().Msg("MQTT repository or DB pool is nil!")
			return
		}

		var intercomData models.CreateIntercomCredentials

		if err := json.Unmarshal(payload, &intercomData); err != nil {
			s.logger.Error().
				Err(err).
				Str("payload", string(payload)).
				Msg("Failed to parse intercom creation message from client")
			return
		}
		s.logger.Info().
			Str("mac", intercomData.MAC).
			Str("address", intercomData.Address).
			Int("apartments", intercomData.NumberOfApartments).
			Msg("Received new intercom creation request from client.")

		intercomID, isNew, err := s.mqqtServerRepository.CreateNewIntercom(intercomData, s.logger)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("mac", intercomData.MAC).
				Msg("Failed to create intercom in database")
			errorResponse := map[string]interface{}{
				"success": false,
			}
			errTopic := fmt.Sprintf("intercom/fromserver/%s", intercomData.MAC)
			if err := s.Publish(errTopic, errorResponse); err != nil {
				s.logger.Error().Err(err).Msg("Failed to send error response")
			}
			return
		}
		responseTopic := fmt.Sprintf("intercom/fromserver/%s", intercomData.MAC)
		response := map[string]interface{}{
			"success": true,
			"id":      intercomID,
			"is_new":  isNew,
			"mac":     intercomData.MAC,
		}
		if err := s.Publish(responseTopic, response); err != nil {
			s.logger.Error().
				Err(err).
				Str("topic", responseTopic).
				Msg("Failed to send intercom creation response")
			return
		}
		s.logger.Info().
			Str("topic", responseTopic).
			Int("id", intercomID).
			Bool("is_new", isNew).
			Msg("Successfully sent intercom creation response")

	})

	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to subscribe to intercom creation topic")
		return
	}

	<-ctx.Done()
	s.logger.Info().Msg("Stopping MQTT listener for intercom creations")
}

func (s *Server) Publish(topic string, payload interface{}) error {
	var data []byte
	switch v := payload.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		var err error
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}
	token := s.client.Publish(topic, byte(s.config.QOSLevel), false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
