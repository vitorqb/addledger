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

// Unique returns unique items from a list. It retains the order.
func Unique[T comparable](a []T) []T {
	out := []T{}
	seen := map[T]int{}
	for _, e := range a {
		if _, ok := seen[e]; ok {
			continue
		}
		seen[e] = 1
		out = append(out, e)
	}
	return out
}

// Reverse reverses an array in place.
func Reverse[T interface{}](a []T) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
