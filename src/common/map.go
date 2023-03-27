package common

func Values[K comparable, V any](input map[K]V) []V {
	values := make([]V, 0, len(input))
	for _, val := range input {
		values = append(values, val)
	}
	return values
}
