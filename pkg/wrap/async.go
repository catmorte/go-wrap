package wrap

import "sync"

type asyncOut[T any] struct {
	result Out[T]

	sync.Mutex

	resultCh chan Out[T]
}

func (r *asyncOut[T]) waitResult() Out[T] {
	res, ok := <-r.resultCh
	if ok {
		r.result = res
		close(r.resultCh)
	}
	return r.result
}

// ErrorOrNil implements Result.
func (r *asyncOut[T]) ErrorOrNil() error {
	return r.waitResult().ErrorOrNil()
}

// Flat implements Result.
func (r *asyncOut[T]) Flat(onOK func(T), onError func(error)) Out[T] {
	return r.waitResult().Flat(onOK, onError)
}

// GetOrDefault implements Result.
func (r *asyncOut[T]) GetOrDefault(defaultValue T) T {
	return r.waitResult().GetOrDefault(defaultValue)
}

// GetOrNil implements Result.
func (r *asyncOut[T]) GetOrNil() *T {
	return r.waitResult().GetOrNil()
}

// IfError implements Result.
func (r *asyncOut[T]) IfError(onError func(error)) Out[T] {
	return r.waitResult().IfError(onError)
}

// IfOK implements Result.
func (r *asyncOut[T]) IfOK(onOK func(T)) Out[T] {
	return r.waitResult().IfOK(onOK)
}

// IsError implements Result.
func (r *asyncOut[T]) IsError() bool {
	return r.waitResult().IsError()
}

// IsOK implements Result.
func (r *asyncOut[T]) IsOK() bool {
	return r.waitResult().IsOK()
}

func Async[T any](fn func() Out[T]) Out[T] {
	resultCh := make(chan Out[T])
	go func() {
		resultCh <- fn()
	}()
	return &asyncOut[T]{resultCh: resultCh}
}

// Unwrap implements Result.
func (r *asyncOut[T]) Unwrap() (T, error) {
	return r.waitResult().Unwrap()
}
