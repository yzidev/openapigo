package goas

// MergeOptions merges a base option slice with additional options,
// returning a new slice (doesn't mutate base).
func MergeOptions(base []HandlerOption, add ...HandlerOption) []HandlerOption {
	out := make([]HandlerOption, 0, len(base)+len(add))
	out = append(out, base...)
	out = append(out, add...)
	return out
}

// MergeOptionSlices merges multiple slices, returning a new slice.
func MergeOptionSlices(slices ...[]HandlerOption) []HandlerOption {
	n := 0
	for _, s := range slices {
		n += len(s)
	}
	out := make([]HandlerOption, 0, n)
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}
