package entities

func FilterEntitiesWithCommands(entityList []Entity) []EntityWithCommand {
	var result []EntityWithCommand
	for _, item := range entityList {
		if v, ok := item.(EntityWithCommand); ok {
			result = append(result, v)
		}
	}

	return result
}
