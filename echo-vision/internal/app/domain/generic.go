package domain

func toStringValues[T ~string](values []T) []string {
	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = string(v)
	}
	return stringValues
}

func isIn[T comparable](s T, values []T) bool {
	for _, v := range values {
		if s == v {
			return true
		}
	}
	return false
}
