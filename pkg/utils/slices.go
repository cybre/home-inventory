package utils

func Map[S ~[]E, E, R any](s S, f func(uint, E) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(uint(i), v)
	}

	return result
}

func MapWithError[S ~[]E, E, R any](s S, f func(uint, E) (R, error)) ([]R, error) {
	result := make([]R, len(s))
	for i, v := range s {
		res, err := f(uint(i), v)
		if err != nil {
			return nil, err
		}

		result[i] = res
	}

	return result, nil
}
