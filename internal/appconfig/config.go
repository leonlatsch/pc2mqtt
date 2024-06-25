package appconfig

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/leonlatsch/windows-hass-bridge/internal/system"
)

const CONFIG_FILE_LOCATION = "config.json"
const CONFIG_FILE_MODE = 0644

var localConfig *AppConfig = nil

func RequireConfig() *AppConfig {
	if localConfig == nil {
		panic("Required app config is nil")
	}

	return localConfig
}

func ValidateOrCreateConfig() error {
	if !configExists() {
		if err := createEmptyConfig(); err != nil {
			return err
		}
		return errors.New("Config does not exist. Created initial config")
	}

	return nil
}

func createEmptyConfig() error {
	newEmptyConfig := AppConfig{
		DeviceId:   uuid.New().String(),
		DeviceName: system.Hostname(),
		Mqtt: MqttAppConfig{
			Host:     "<YOUR MQTT HOST>",
			Port:     1883,
			Username: "<MQTT USER>",
			Password: "<MQTT PASSWORD>",
		},
		DebugMode: false,
	}

	if err := SaveConfig(newEmptyConfig); err != nil {
		return err
	}

	return nil
}

func configExists() bool {
	_, err := os.Stat(CONFIG_FILE_LOCATION)

	return !os.IsNotExist(err)
}

func SaveConfig(conf AppConfig) error {
	confJson, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(CONFIG_FILE_LOCATION, confJson, CONFIG_FILE_MODE); err != nil {
		return err
	}

	localConfig = &conf
	return nil
}

func LoadConfig() error {
	var conf AppConfig
	buf, err := os.ReadFile(CONFIG_FILE_LOCATION)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, &conf); err != nil {
		return err
	}

	localConfig = &conf
	return nil
}
