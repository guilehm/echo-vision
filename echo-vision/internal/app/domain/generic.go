package domain

func toStringValues[T ~string](values []T) []string {
	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = string(v)
	}
	return stringValues
}
