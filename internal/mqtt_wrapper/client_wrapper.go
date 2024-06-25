package mqtt_wrapper

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-mqtt/mqtt"
)

type MqttClientConfig struct {
	ClientId  string
	Uri       string
	Username  string
	Password  string
	SubTopics []string
	Will      Message
}

type MqttClientWrapper struct {
	Config      *MqttClientConfig
	InnerClient *mqtt.Client
}

type Message struct {
	Topic   string
	Message string
}

func (client *MqttClientWrapper) Subscribe(topics ...string) chan Message {
	var messageChan = make(chan Message)

	// Send sub request
	go func() {
		log.Println("Subsribing to " + fmt.Sprint(len(topics)) + " topics")

		err := client.InnerClient.Subscribe(nil, topics...)
		if err != nil {
			panic(err)
		}
	}()

	// Read messages
	go func() {
		for {
			message, topic, err := client.InnerClient.ReadSlices()

			switch {
			case err == nil:
				messageObj := Message{
					Topic:   string(topic),
					Message: string(message),
				}

				messageChan <- messageObj
			case errors.Is(err, mqtt.ErrClosed):
				log.Println(err)
				return
			case mqtt.IsConnectionRefused(err):
				log.Println(err)
				time.Sleep(15 * time.Second)
			default:
				log.Println(err)
				time.Sleep(2 * time.Second)
			}
		}
	}()

	return messageChan
}

func (client *MqttClientWrapper) Publish(topic string, message []byte) error {
	return client.InnerClient.Publish(nil, message, topic)
}

func CreateClientWrapper(config *MqttClientConfig) *MqttClientWrapper {
	innerClient, err := mqtt.VolatileSession(config.ClientId, &mqtt.Config{
		Dialer:       mqtt.NewDialer("tcp", config.Uri),
		PauseTimeout: 4 * time.Second,
		UserName:     config.Username,
		Password:     []byte(config.Password),
		Will: struct {
			Topic       string
			Message     []byte
			Retain      bool
			AtLeastOnce bool
			ExactlyOnce bool
		}{
			Topic:   config.Will.Topic,
			Message: []byte(config.Will.Message),
			Retain:  true,
		},
	})

	if err != nil {
		panic(err)
	}

	log.Println("Connecting to " + config.Uri)

	return &MqttClientWrapper{
		Config:      config,
		InnerClient: innerClient,
	}
}
