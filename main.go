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

	client := createClient()

	for {
		ctx, cancel := context.WithCancel(context.Background())

		entityList := entities.GetEntities()
		entitiesWithCommands := entities.FilterEntitiesWithCommands(entityList)

		go logWhenClientOnline(client)

		go func() {
			// Wait with availability until config was sent
			publishAutoDiscoveryConfigs(ctx, cancel, client, entityList)
			publishAvailability(ctx, cancel, client, entityList)
		}()

		go publishSensorStates(ctx, cancel, client, entityList)
		go subToCmdTopics(ctx, cancel, client, entitiesWithCommands)

		go readMessages(cancel, client, entitiesWithCommands)

		<-ctx.Done()
		log.Println("Connection lost. Reconnecting in 5 seconds...")

		time.Sleep(5 * time.Second)
		cancel()
	}
}

func logWhenClientOnline(client *mqtt.Client) {
	appConf := appconfig.RequireConfig()
	<-client.Online()
	log.Printf("Conencted to %q", appConf.Mqtt.Host)
}

func publishAutoDiscoveryConfigs(ctx context.Context, cancel func(), client *mqtt.Client, entityList []entities.Entity) {
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

	log.Println("Published auto discovery messages successfully")
}

func publishAvailability(ctx context.Context, cancel func(), client *mqtt.Client, entityList []entities.Entity) {
	for _, ety := range entityList {
		availability := ety.GetDiscoveryConfig().Availability
		payload := []byte(availability.PayloadAvailable)
		if err := client.PublishRetained(ctx.Done(), payload, availability.Topic); err != nil {
			log.Println(err)
			cancel()
		}
	}

	log.Println("Published availability messages successfully")
}

func publishSensorStates(ctx context.Context, cancel func(), client *mqtt.Client, entityList []entities.Entity) {
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

	debugLog("Published sensor state successfully")
}

func subToCmdTopics(ctx context.Context, cancel func(), client *mqtt.Client, entitiesWithCommands []entities.EntityWithCommand) {
	cmdTopics := make([]string, 0, len(entitiesWithCommands))
	for _, ety := range entitiesWithCommands {
		cmdTopics = append(cmdTopics, ety.GetDiscoveryConfig().CommandTopic)
	}

	if err := client.Subscribe(ctx.Done(), cmdTopics...); err != nil {
		log.Println(err)
		cancel()
	}

	debugLog("Subscribed to command topic")
}

func readMessages(cancel func(), client *mqtt.Client, entitiesWithCommands []entities.EntityWithCommand) {
	for {
		message, topic, err := client.ReadSlices()

		switch {
		case err == nil:
			top := string(topic)
			mes := string(message)

			debugLog("DEBUG: REC: " + top + " : " + mes)

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

func debugLog(message string) {
	if appconfig.RequireConfig().DebugMode {
		log.Println(message)
	}
}
