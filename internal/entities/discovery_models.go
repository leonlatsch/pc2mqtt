package entities

type Config struct {
	Device       Device       `json:"device"`
	Availability Availability `json:"availability"`
	CommandTopic string       `json:"command_topic"`
	Name         string       `json:"name"`
	Icon         string       `json:"icon"`
	ObjectId     string       `json:"object_id"`
	StateTopic   string       `json:"state_topic"`
	PayloadOn    string       `json:"payload_on"`
	PayloadOff   string       `json:"payload_off"`
	Schema       string       `json:"schema"`
}

type Device struct {
	Identifiers  string `json:"identifiers"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Name         string `json:"name"`
}

type Availability struct {
	Topic               string `json:"topic"`
	PayloadAvailable    string `json:"payload_available"`
	PayloadNotAvailable string `json:"payload_not_available"`
}
