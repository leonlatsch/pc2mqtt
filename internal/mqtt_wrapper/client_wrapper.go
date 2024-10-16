package mqtt_wrapper

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-mqtt/mqtt"
	"github.com/leonlatsch/pc2mqtt/internal/appconfig"
	"github.com/leonlatsch/pc2mqtt/internal/entities"
)

type MqttClientWrapper struct {
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
		if appconfig.RequireConfig().DebugMode {
			for _, topic := range topics {
				log.Println("DEBUG: Subscribing to " + topic)
			}
		}

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
				log.Fatalln(err)
			case mqtt.IsConnectionRefused(err):
				log.Fatalln(err)
			default:
				log.Fatalln(err)
			}
		}
	}()

	return messageChan
}

func (client *MqttClientWrapper) Publish(topic string, message []byte) error {
	if appconfig.RequireConfig().DebugMode {
		log.Println("DEBUG: PUB: " + topic)
	}
	return client.InnerClient.Publish(nil, message, topic)
}

func CreateClientWrapper() *MqttClientWrapper {
	appConf := appconfig.RequireConfig()
	clientId := "pc2mqtt-" + appConf.DeviceName
	uri := fmt.Sprintf("%v:%v", appConf.Mqtt.Host, appConf.Mqtt.Port)

	innerClient, err := mqtt.VolatileSession(clientId, &mqtt.Config{
		Dialer:       mqtt.NewDialer("tcp", uri),
		PauseTimeout: 4 * time.Second,
		UserName:     appConf.Mqtt.Username,
		Password:     []byte(appConf.Mqtt.Password),
		Will: struct {
			Topic       string
			Message     []byte
			Retain      bool
			AtLeastOnce bool
			ExactlyOnce bool
		}{
			Topic:   entities.GetDeviceAvailability().Topic,
			Message: []byte(entities.GetDeviceAvailability().PayloadNotAvailable),
			Retain:  true,
		},
	})

	if err != nil {
		panic(err)
	}

	log.Println("Connecting to " + uri)

	return &MqttClientWrapper{
		InnerClient: innerClient,
	}
}
