package common

func Keys[K comparable, V any](input map[K]V) []K {
	keys := make([]K, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	return keys
}

func Values[K comparable, V any](input map[K]V) []V {
	values := make([]V, 0, len(input))
	for _, val := range input {
		values = append(values, val)
	}
	return values
}

func FilterMap[K comparable, V any](input map[K]V, predicate func(K, V) bool) map[K]V {
	filteredMap := make(map[K]V)

	for key, value := range input {
		if predicate(key, value) {
			filteredMap[key] = value
		}
	}

	return filteredMap
}
