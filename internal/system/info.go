package system

import (
	"os"
)

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "DEFAULT_SYSTEM_NAME"
	}

	return hostname
}
