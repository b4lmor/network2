package utils

func RemoveElement[T comparable](slice []T, value T) []T {
	var result []T
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

func FindByPredicate[T comparable](slice []*T, predicate func(*T) bool) *T {
	for _, v := range slice {
		if predicate(v) {
			return v
		}
	}
	return nil
}
