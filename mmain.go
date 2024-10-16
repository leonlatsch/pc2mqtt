package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-mqtt/mqtt"
	"github.com/leonlatsch/pc2mqtt/internal/appconfig"
	"github.com/leonlatsch/pc2mqtt/internal/entities"
)

func mainn() {
	log.Println("Starting application")

	if err := appconfig.LoadConfig(); err != nil {
		log.Fatalln(err)
	}
	appConf := appconfig.RequireConfig()

	client := createClient()

	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		entityList := entities.GetEntities()
		entitiesWithCommands := entities.FilterEntitiesWithCommands(entityList)

		// Publish auto discovery
		go func() {
			for _, ety := range entityList {
				configJson, err := json.Marshal(ety.GetDiscoveryConfig())
				if err != nil {
					log.Println(err)
					continue
				}

				if err := client.Publish(ctx.Done(), []byte(configJson), ety.GetDiscoveryTopic()); err != nil {
					log.Println(err)
				}

			}
		}()
		// Publish availability

		go func() {
			availability := entities.GetFixAvailability()
			payload := []byte(availability.PayloadAvailable)
			client.Publish(ctx.Done(), payload, availability.Topic)
		}()

		// Publish sensor state

		cmdTopics := make([]string, 0, len(entitiesWithCommands))
		for _, ety := range entitiesWithCommands {
			cmdTopics = append(cmdTopics, ety.GetDiscoveryConfig().CommandTopic)
		}

		go func() {
			client.Subscribe(ctx.Done(), cmdTopics...)
		}()

		go func() {
			for {
				message, topic, err := client.ReadSlices()

				switch {
				case err == nil:
					// handles message
					if appConf.DebugMode {
						top := string(topic)
						mes := string(message)
						log.Println(top + " : " + mes)
					}
				default:
					// Cancel context and return go func
					cancel()
					return
				}
			}
		}()

		<-ctx.Done()
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
			Topic:   entities.GetFixAvailability().Topic,
			Message: []byte(entities.GetFixAvailability().PayloadNotAvailable),
			Retain:  true,
		},
	})

	if err != nil {
		panic(err)
	}

	return client
}
