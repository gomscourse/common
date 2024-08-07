package tools

func MapSlice[T any, R any](in []T, mapFunc func(T) R) []R {
	out := make([]R, len(in))
	for i := range in {
		out[i] = mapFunc(in[i])
	}
	return out
}

func FilterSlice[T any](in []T, filterFunc func(T) bool) []T {
	var out []T
	for i := range in {
		if filterFunc(in[i]) {
			out = append(out, in[i])
		}
	}
	return out
}
