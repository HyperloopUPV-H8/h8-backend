package common

func Remove[T any](slice []T, i int) (after []T, removed T) {
	removed = slice[i]
	after = append(slice[:i], slice[i+1:]...)
	return after, removed
}
