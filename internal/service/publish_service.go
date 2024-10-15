package service

import (
	"encoding/json"
	"log"

	"github.com/leonlatsch/pc2mqtt/internal/entities"
	"github.com/leonlatsch/pc2mqtt/internal/mqtt_wrapper"
)

type MqttPublisherService struct {
	Client *mqtt_wrapper.MqttClientWrapper
}

func (service *MqttPublisherService) PublishOnStartup(entityList []entities.Entity) {
	service.PublishAutoDiscoveryMessages(entityList)
	service.PublishAvailability()

	var sensors []entities.BinarySensor
	for _, entity := range entityList {
		switch v := entity.(type) {
		case entities.BinarySensor:
			sensors = append(sensors, v)
		}
	}
	service.PublishSensorStates(sensors)
}

func (service *MqttPublisherService) PublishAutoDiscoveryMessages(entityList []entities.Entity) {
	for _, entity := range entityList {
		configJson, err := json.Marshal(entity.GetDiscoveryConfig())
		if err != nil {
			log.Println(err)
			continue
		}

		// Fire and forget
		go service.Client.Publish(entity.GetDiscoveryTopic(), configJson)
	}
}

func (service *MqttPublisherService) PublishAvailability() {
	topic := entities.GetFixAvailability().Topic
	payload := []byte(entities.GetFixAvailability().PayloadAvailable)
	go service.Client.Publish(topic, payload)
}

func (service *MqttPublisherService) PublishSensorStates(sensors []entities.BinarySensor) {
	for _, sensor := range sensors {
		go service.Client.Publish(sensor.DiscoveryConfig.StateTopic, []byte("ON"))
	}
}
