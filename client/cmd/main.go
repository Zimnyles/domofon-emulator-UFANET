package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Domofon struct {
	MAC      string
	DoorOpen bool
	Active   bool
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + os.Getenv("MQTT_HOST") + ":1883")
	opts.SetClientID("domofon-" + generateMAC())
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	domofon := Domofon{
		MAC:      generateMAC(),
		DoorOpen: false,
		Active:   true,
	}

	// Подписка на команды открытия
	client.Subscribe("domofon/"+domofon.MAC+"/open", 0, func(c mqtt.Client, m mqtt.Message) {
		log.Println("Door opened by server command")
		domofon.DoorOpen = true
		client.Publish("domofon/"+domofon.MAC+"/status", 0, false, createStatus(domofon))
	})

	// Периодическая отправка статуса
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			client.Publish("domofon/status", 0, false, createStatus(domofon))
		}
	}()

	// Имитация событий
	for {
		time.Sleep(time.Duration(rand.Intn(120)+30) * time.Second)
		eventType := rand.Intn(2)

		switch eventType {
		case 0: // Звонок
			client.Publish("domofon/events", 0, false, createEvent(domofon, "call"))
		case 1: // Открытие ключом
			client.Publish("domofon/events", 0, false, createEvent(domofon, "key_open"))
			domofon.DoorOpen = true
			client.Publish("domofon/status", 0, false, createStatus(domofon))
		}
	}
}

func createStatus(d Domofon) string {
	return `{"mac":"` + d.MAC + `","door_open":` + strconv.FormatBool(d.DoorOpen) + `,"active":true}`
}

func createEvent(d Domofon, eventType string) string {
	apt := rand.Intn(100) + 1
	return `{"mac":"` + d.MAC + `","type":"` + eventType + `","apartment":` + strconv.Itoa(apt) + `}`
}

func generateMAC() string {
	mac := make([]byte, 6)
	rand.Read(mac)
	return "02" + // Локально администрируемый адрес
		":" + strconv.FormatInt(int64(mac[0]), 16) +
		":" + strconv.FormatInt(int64(mac[1]), 16) +
		":" + strconv.FormatInt(int64(mac[2]), 16) +
		":" + strconv.FormatInt(int64(mac[3]), 16) +
		":" + strconv.FormatInt(int64(mac[4]), 16)
}
