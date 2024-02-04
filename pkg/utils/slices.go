package utils

func Map[S ~[]E, E, R any](s S, f func(uint, E) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(uint(i), v)
	}

	return result
}
