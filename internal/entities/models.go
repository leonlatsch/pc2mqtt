package entities

type Entity interface {
	GetDiscoveryTopic() string
	GetDiscoveryConfig() *DiscoveryConfig
}

type EntityWithCommand interface {
	Entity
	QueueAction()
}

// https://www.home-assistant.io/integrations/binary_sensor.mqtt
type BinarySensor struct {
	DiscoveryTopic  string
	DiscoveryConfig *DiscoveryConfig
}

func (sensor BinarySensor) GetDiscoveryTopic() string {
	return sensor.DiscoveryTopic
}

func (sensor BinarySensor) GetDiscoveryConfig() *DiscoveryConfig {
	return sensor.DiscoveryConfig
}

type Button struct {
	DiscoveryTopic  string
	DiscoveryConfig *DiscoveryConfig
	Action          func()
}

func (button Button) GetDiscoveryTopic() string {
	return button.DiscoveryTopic
}

func (button Button) GetDiscoveryConfig() *DiscoveryConfig {
	return button.DiscoveryConfig
}

func (button Button) QueueAction() {
	go button.Action()
}

// Disco
