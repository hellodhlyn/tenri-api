package utils

func MapSlice[S any, D any](src []S, fn func(S) D) []D {
	dst := make([]D, len(src))
	for i, each := range src {
		dst[i] = fn(each)
	}
	return dst
}
