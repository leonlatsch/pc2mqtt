package appconfig

type MqttAppConfig struct {
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	AutoDiscoveryPrefix string `json:"auto_discovery_prefix"`
}

type AppConfig struct {
	DeviceId   string        `json:"device_id"`
	DeviceName string        `json:"device_name"`
	Mqtt       MqttAppConfig `json:"mqtt"`
	DebugMode  bool          `json:"debug_mode"`
}
