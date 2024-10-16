package appconfig

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/leonlatsch/pc2mqtt/internal/system"
)

const configFileName = "config.json"
const configFileMode = 0644

var localConfig *AppConfig = nil

func RequireConfig() *AppConfig {
	if localConfig == nil {
		panic("Required app config is nil")
	}

	return localConfig
}

func createEmptyConfig() error {
	newEmptyConfig := AppConfig{
		DeviceId:   uuid.New().String(),
		DeviceName: system.Hostname(),
		Mqtt: MqttAppConfig{
			Host:                "YOUR MQTT HOST",
			Port:                1883,
			Username:            "MQTT USER",
			Password:            "MQTT PASSWORD",
			AutoDiscoveryPrefix: "homeassistant",
		},
		DebugMode: false,
	}

	if err := SaveConfig(newEmptyConfig); err != nil {
		return err
	}

	return nil
}

func configExists() bool {
	_, err := os.Stat(configFileName)
	return !os.IsNotExist(err)
}

func SaveConfig(conf AppConfig) error {
	confJson, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFileName, confJson, configFileMode); err != nil {
		return err
	}

	localConfig = &conf
	return nil
}

func LoadConfig() error {
	if !configExists() {
		if err := createEmptyConfig(); err != nil {
			return err
		}
		return errors.New("Config does not exist. Created initial config")
	}

	var conf AppConfig
	buf, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, &conf); err != nil {
		return err
	}

	localConfig = &conf
	return nil
}
