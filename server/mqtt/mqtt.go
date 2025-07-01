package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Connect(url string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID("mqtt-golang")
	client := mqtt.NewClient(opts)
	fmt.Println("1")
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}
