package wrap

type err[T any] struct {
	err error
}

func (e err[T]) ErrorOrNil() error {
	return e.err
}

func (e err[T]) Flat(_ func(T), onError func(error)) Output[T] {
	onError(e.err)
	return e
}

func (err[T]) GetOrDefault(defaultValue T) T {
	return defaultValue
}

func (err[T]) GetOrNil() *T {
	return nil
}

func (err[T]) IsError() bool {
	return true
}

func (err[T]) IsOK() bool {
	return false
}

func (f err[T]) IfError(onError func(error)) Output[T] {
	onError(f.err)
	return f
}

func (e err[T]) IfOK(onOK func(T)) Output[T] {
	return e
}

func Err[T any](e error) Output[T] {
	output := new(err[T])
	output.err = e

	return *output
}

func (e err[T]) Unwrap() (T, error) {
	var zeroVal T
	return zeroVal, e.err
}
