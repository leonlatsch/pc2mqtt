package main

import (
	"fmt"
	"log"

	"github.com/leonlatsch/pc2mqtt/internal/appconfig"
	"github.com/leonlatsch/pc2mqtt/internal/entities"
	"github.com/leonlatsch/pc2mqtt/internal/ext"
	"github.com/leonlatsch/pc2mqtt/internal/mqtt_wrapper"
	"github.com/leonlatsch/pc2mqtt/internal/service"
)

func main() {
	log.Println("Starting application")

	if err := appconfig.LoadConfig(); err != nil {
		log.Println(err)
		return
	}
	appConfig := appconfig.RequireConfig()

	clientWrapper := mqtt_wrapper.CreateClientWrapper()
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

		if appConfig.DebugMode {
			log.Println("DEBUG: REC: " + fmt.Sprintf("%v -> %v", message.Topic, message.Message))
		}

		for _, entity := range entitiesWithCommands {
			if entity.GetDiscoveryConfig().CommandTopic == message.Topic {
				entity.QueueAction()
			}
		}
	}
}
