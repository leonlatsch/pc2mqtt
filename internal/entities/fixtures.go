package entities

import (
	"log"
	"runtime"

	"github.com/leonlatsch/pc2mqtt/internal/appconfig"
	"github.com/leonlatsch/pc2mqtt/internal/system"
)

func GetEntities() []Entity {
	appConf := appconfig.RequireConfig()
	entityList := []Entity{
		BinarySensor{
			DiscoveryTopic: "homeassistant/binary_sensor/" + appConf.DeviceId + "/" + appConf.DeviceName + "_sensor_power/config",
			DiscoveryConfig: &Config{
				Device:       GetFixDevice(),
				Availability: GetFixAvailability(),
				ObjectId:     appConf.DeviceName + "_sensor_power",
				UniqueId:     appConf.DeviceName + "_sensor_power",
				Name:         "Power",
				Icon:         "mdi:power",
				StateTopic:   appConf.DeviceName + "/binary_sensor/power/state",
				Qos:          1,
			},
		},
		Button{
			Action: func() {
				cmd, err := system.GetShutdownCommand()
				if err != nil {
					log.Println(err)
					return
				}

				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
			},
			DiscoveryTopic: "homeassistant/button/" + appConf.DeviceId + "/" + appConf.DeviceName + "_button_shutdown/config",
			DiscoveryConfig: &Config{
				Device:       GetFixDevice(),
				Availability: GetFixAvailability(),
				ObjectId:     appConf.DeviceName + "_button_shutdown",
				UniqueId:     appConf.DeviceName + "_button_shutdown",
				Name:         "Shutdown",
				Icon:         "mdi:power",
				StateTopic:   appConf.DeviceName + "/button/shutdown/state",
				CommandTopic: appConf.DeviceName + "/button/shutdown/command",
				Qos:          1,
			},
		},
		Button{
			Action: func() {
				cmd, err := system.GetRebootCommand()
				if err != nil {
					log.Println(err)
					return
				}

				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
			},
			DiscoveryTopic: "homeassistant/button/" + appConf.DeviceId + "/" + appConf.DeviceName + "_button_reboot/config",
			DiscoveryConfig: &Config{
				Device:       GetFixDevice(),
				Availability: GetFixAvailability(),
				ObjectId:     appConf.DeviceName + "_button_reboot",
				UniqueId:     appConf.DeviceName + "_button_reboot",
				Name:         "Reboot",
				Icon:         "mdi:restart",
				StateTopic:   appConf.DeviceName + "/button/reboot/state",
				CommandTopic: appConf.DeviceName + "/button/reboot/command",
				Qos:          1,
			},
		},
	}

	if appConf.DebugMode {
		entityList = append(entityList,
			Button{
				Action: func() {
					log.Println("Test button pressed")
				},
				DiscoveryTopic: "homeassistant/button/" + appConf.DeviceId + "/" + appConf.DeviceName + "_button_test/config",
				DiscoveryConfig: &Config{
					Device:       GetFixDevice(),
					Availability: GetFixAvailability(),
					ObjectId:     appConf.DeviceName + "_button_test",
					UniqueId:     appConf.DeviceName + "_button_test",
					Name:         "Test",
					Icon:         "mdi:test-tube",
					StateTopic:   appConf.DeviceName + "/button/test/state",
					CommandTopic: appConf.DeviceName + "/button/test/command",
					Qos:          1,
				},
			},
		)
	}

	return entityList
}

func GetFixAvailability() Availability {
	appConf := appconfig.RequireConfig()
	return Availability{
		Topic:               appConf.DeviceName + "/state",
		PayloadAvailable:    "online",
		PayloadNotAvailable: "offline",
	}

}
func GetFixDevice() Device {
	appConf := appconfig.RequireConfig()
	return Device{
		Identifiers:  appConf.DeviceId,
		Manufacturer: runtime.GOOS + "/" + runtime.GOARCH,
		Model:        appConf.DeviceName,
		Name:         appConf.DeviceName,
	}
}
