package util

func Contains[T comparable](xs []T, v T) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func ContainsAny[T comparable](xs []T, vs []T) bool {
	for _, x := range xs {
		for _, v := range vs {
			if x == v {
				return true
			}
		}
	}
	return false
}
