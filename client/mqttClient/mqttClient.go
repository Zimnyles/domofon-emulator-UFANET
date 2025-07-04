package mqttclient

import (
	"domofonEmulator/config"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type Client struct {
	client mqtt.Client
	logger *zerolog.Logger
	config config.MQTTConfig
}

func Connect(mqqtConfig config.MQTTConfig, logger *zerolog.Logger) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqqtConfig.Broker)
	opts.SetClientID("client")
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
	return &Client{
		client: client,
		logger: logger,
	}, nil
}

func (c *Client) Subscribe(topic string, handler func(payload []byte)) error {
	token := c.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) Publish(topic string, payload interface{}) error {
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
	token := c.client.Publish(topic, byte(c.config.QOSLevel), false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
	c.logger.Info().Msg("MQTT client disconnected")
}
