package utils

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func MapToSlice[T any](m map[string]interface{}) []T {
	slice := make([]T, 0, len(m))
	for _, v := range m {
		slice = append(slice, v.(T))
	}
	return slice
}
