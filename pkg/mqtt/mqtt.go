package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func Connect(url string, logger *zerolog.Logger) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID("mqtt-golang")
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MQTT broker")
		return nil

	} else {
		logger.Info().Msg("Client is connected to MQTT broker")
	}
	return client
}

