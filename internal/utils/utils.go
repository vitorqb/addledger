package utils

// StopSplitArray signals that the SplitArray should stop.
type StopSplitArray struct{}

func (*StopSplitArray) Error() string { return "Stop!" }

// SplitArray splits an array into sub-arrays
func SplitArray[T interface{}](size int, a []T) (next func() ([]T, error), err error) {
	from := 0
	to := size
	next = func() ([]T, error) {
		if to <= len(a) {
			out := a[from:to]
			from += size
			to += size
			return out, nil
		}
		return nil, &StopSplitArray{}
	}
	return next, nil
}

// RemoveIndex removes an entry from an array
func RemoveIndex[T interface{}](i int, a []T) []T {
	return append(a[:i], a[i+1:]...)
}
