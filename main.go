package main

import (
	"fmt"
	"log"

	"github.com/leonlatsch/windows-hass-bridge/internal/appconfig"
	"github.com/leonlatsch/windows-hass-bridge/internal/entities"
	"github.com/leonlatsch/windows-hass-bridge/internal/ext"
	"github.com/leonlatsch/windows-hass-bridge/internal/mqtt_wrapper"
	"github.com/leonlatsch/windows-hass-bridge/internal/service"
)

func main() {
	log.Println("Running application")

	if err := appconfig.ValidateOrCreateConfig(); err != nil {
		log.Println("Config not valid. Empty config was created")
		return
	}

	if err := appconfig.LoadConfig(); err != nil {
		log.Panicln(err)
	}
	appConfig := appconfig.RequireConfig()

	// Test mqtt library
	clientWrapper := mqtt_wrapper.CreateClientWrapper(&mqtt_wrapper.MqttClientConfig{
		ClientId: "hass-bridge-" + appConfig.DeviceName,
		Uri:      fmt.Sprintf("%v:%v", appConfig.Mqtt.Host, appConfig.Mqtt.Port),
		Username: appConfig.Mqtt.Username,
		Password: appConfig.Mqtt.Password,
		Will: mqtt_wrapper.Message{
			Topic:   entities.GetFixAvailability().Topic,
			Message: entities.GetFixAvailability().PayloadNotAvailable,
		},
	})
	publishService := service.MqttPublisherService{
		Client: clientWrapper,
	}

	entityList := entities.GetEntities()
	entitiesWithCommands := ext.FilterEntiiesWithCommands(entityList)

	publishService.PublishOnStartup(entityList)

	cmdTopics := make([]string, 0, len(entitiesWithCommands))
	for _, k := range entitiesWithCommands {
		cmdTopics = append(cmdTopics, k.GetDiscoveryConfig().CommandTopic)
	}

	messagesChan := clientWrapper.Subscribe(cmdTopics...)

	for {
		var message = <-messagesChan
		for _, entity := range entitiesWithCommands {
			if entity.GetDiscoveryConfig().CommandTopic == message.Topic {
				entity.QueueAction()
			}
		}
	}
}
