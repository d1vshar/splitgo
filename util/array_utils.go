package util

func Transform[T any, U any](arr []T, f func(T) U) []U {
	result := make([]U, len(arr))
	for i, v := range arr {
		result[i] = f(v)
	}
	return result
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Any[T any](ss []T, test func(T) bool) bool {
	for _, s := range ss {
		if test(s) {
			return true
		}
	}
	return false
}
