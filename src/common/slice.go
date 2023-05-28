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

func Map[T any, U any](items []T, mapFn func(item T) U) []U {
	result := make([]U, len(items))

	for index, item := range items {
		result[index] = mapFn(item)
	}

	return result
}

func Every[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}

	return true
}

func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for index, item := range slice {
		if predicate(item) {
			return index
		}
	}

	return -1
}
