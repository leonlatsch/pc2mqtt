package entities

type Entity interface {
	GetDiscoveryTopic() string
	GetDiscoveryConfig() *Config
}

type EntityWithCommand interface {
	Entity
	QueueAction()
}

type BinarySensor struct {
	DiscoveryTopic  string
	DiscoveryConfig *Config
}

func (sensor BinarySensor) GetDiscoveryTopic() string {
	return sensor.DiscoveryTopic
}

func (sensor BinarySensor) GetDiscoveryConfig() *Config {
	return sensor.DiscoveryConfig
}

type Button struct {
	DiscoveryTopic  string
	DiscoveryConfig *Config
	Action          func()
}

func (button Button) GetDiscoveryTopic() string {
	return button.DiscoveryTopic
}

func (button Button) GetDiscoveryConfig() *Config {
	return button.DiscoveryConfig
}

func (button Button) QueueAction() {
	go button.Action()
}

// Disco
