package ext

func Map[F any, T any](list []F, transform func(F) T) []T {
	newList := make([]T, len(list))

	for i, e := range list {
		newList[i] = transform(e)
	}

	return newList
}
