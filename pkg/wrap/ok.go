package wrap

type ok[T any] struct {
	v T
}

func (ok[T]) ErrorOrNil() error {
	return nil
}

func (s ok[T]) Flat(onOK func(T), onError func(error)) Out[T] {
	onOK(s.v)
	return s
}

func (s ok[T]) GetOrDefault(defaultValue T) T {
	return s.v
}

func (s ok[T]) GetOrNil() *T {
	return &s.v
}

func (ok[T]) IsError() bool {
	return false
}

func (ok[T]) IsOK() bool {
	return true
}

func (s ok[T]) IfError(onError func(error)) Out[T] {
	return s
}

func (s ok[T]) IfOK(onOK func(T)) Out[T] {
	onOK(s.v)
	return s
}

func OK[T any](value T) Out[T] {
	output := new(ok[T])
	output.v = value

	return *output
}

func OKVargs[T any](values ...T) []Out[T] {
	res := make([]Out[T], 0, len(values))
	for _, v := range values {
		res = append(res, OK(v))
	}
	return res
}

func OKSlice[T any](values []T) []Out[T] {
	res := make([]Out[T], 0, len(values))
	for _, v := range values {
		res = append(res, OK(v))
	}
	return res
}

func (s ok[T]) Unwrap() (T, error) {
	return s.v, nil
}
