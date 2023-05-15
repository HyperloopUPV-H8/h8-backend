package common

func Remove[T any](slice []T, i int) (after []T, removed T) {
	removed = slice[i]
	after = append(slice[:i], slice[i+1:]...)
	return after, removed
}

func Filter[T any](items []T, predicate func(item T) bool) []T {
	newSlice := make([]T, 0)

	for _, item := range items {
		if predicate(item) {
			newSlice = append(newSlice, item)
		}
	}

	return newSlice
}
