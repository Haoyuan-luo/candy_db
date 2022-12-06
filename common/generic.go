package common

func GenZero[T any]() (r T) {
	return
}

func GenPtr[T any](t T) *T {
	return &t
}

func GenDirect[T any](t *T) T {
	if t == nil {
		return GenZero[T]()
	}
	return *t
}

func GenUnique[T comparable](t []T) []T {
	m := make(map[T]bool)
	for _, v := range t {
		m[v] = true
	}
	r := make([]T, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func GenEqual[T comparable](a, b T) bool {
	return a == b
}

func GenFilter[T any](t []T, f func(T) bool) []T {
	r := make([]T, 0, len(t))
	for _, v := range t {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func GenFilterNot[T any](t []T, f func(T) bool) []T {
	return GenFilter[T](t, func(v T) bool {
		return !f(v)
	})
}
