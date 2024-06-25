package ext

import "github.com/leonlatsch/windows-hass-bridge/internal/entities"

func FilterEntiiesWithCommands(entityList []entities.Entity) []entities.EntityWithCommand {
	var result []entities.EntityWithCommand
	for _, item := range entityList {
		if v, ok := item.(entities.EntityWithCommand); ok {
			result = append(result, v)
		}
	}

	return result
}
