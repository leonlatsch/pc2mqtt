package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-mqtt/mqtt"
	"github.com/leonlatsch/pc2mqtt/internal/appconfig"
	"github.com/leonlatsch/pc2mqtt/internal/entities"
	"log"
	"time"
)

func main() {
	log.Println("Starting application")

	if err := appconfig.LoadConfig(); err != nil {
		log.Fatalln(err)
	}
	appConf := appconfig.RequireConfig()

	client := createClient()

	for {
		ctx, cancel := context.WithCancel(context.Background())

		entityList := entities.GetEntities()
		entitiesWithCommands := entities.FilterEntitiesWithCommands(entityList)

		go func() {
			// Send auto discovery config for all entities
			for _, ety := range entityList {
				configJson, err := json.Marshal(ety.GetDiscoveryConfig())
				if err != nil {
					log.Println(err)
					continue
				}

				if err := client.PublishRetained(ctx.Done(), configJson, ety.GetDiscoveryTopic()); err != nil {
					log.Println(err)
					cancel()
					return
				}

			}

			// Send payload available for all entities
			for _, ety := range entityList {
				availability := ety.GetDiscoveryConfig().Availability
				payload := []byte(availability.PayloadAvailable)
				if err := client.PublishRetained(ctx.Done(), payload, availability.Topic); err != nil {
					log.Println(err)
					cancel()
				}
			}
		}()

		// Publish sensor state
		go func() {
			var sensors []entities.BinarySensor
			for _, entity := range entityList {
				switch v := entity.(type) {
				case entities.BinarySensor:
					sensors = append(sensors, v)
				}
			}

			for _, sensor := range sensors {
				topic := sensor.GetDiscoveryConfig().StateTopic
				payload := []byte(sensor.DiscoveryConfig.PayloadOn)
				if err := client.PublishRetained(ctx.Done(), payload, topic); err != nil {
					log.Println(err)
					cancel()
					return
				}
			}
		}()

		// Subscribe to command topics
		go func() {
			cmdTopics := make([]string, 0, len(entitiesWithCommands))
			for _, ety := range entitiesWithCommands {
				cmdTopics = append(cmdTopics, ety.GetDiscoveryConfig().CommandTopic)
			}

			if err := client.Subscribe(ctx.Done(), cmdTopics...); err != nil {
				log.Println(err)
				cancel()
			}
		}()

		go func() {
			for {
				message, topic, err := client.ReadSlices()

				switch {
				case err == nil:
					top := string(topic)
					mes := string(message)

					if appConf.DebugMode {
						log.Println("DEBUG: REC: " + top + " : " + mes)
					}

					for _, ety := range entitiesWithCommands {
						if ety.GetDiscoveryConfig().CommandTopic == top {
							ety.QueueAction()
						}
					}

				default:
					cancel()
					return
				}
			}
		}()

		<-ctx.Done()
		log.Println("Connection lost. Reconnecting in 5 seconds...")

		time.Sleep(5 * time.Second)
		cancel()
	}
}

func createClient() *mqtt.Client {
	appConf := appconfig.RequireConfig()
	clientId := "pc2mqtt-" + appConf.DeviceName
	uri := fmt.Sprintf("%v:%v", appConf.Mqtt.Host, appConf.Mqtt.Port)

	client, err := mqtt.VolatileSession(clientId, &mqtt.Config{
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

	return client
}
