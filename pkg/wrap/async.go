package wrap

import "sync"

type asyncOutput[T any] struct {
	result Output[T]

	sync.Mutex

	resultCh chan Output[T]
}

func (r *asyncOutput[T]) waitResult() Output[T] {
	res, ok := <-r.resultCh
	if ok {
		r.result = res
		close(r.resultCh)
	}
	return r.result
}

// ErrorOrNil implements Result.
func (r *asyncOutput[T]) ErrorOrNil() error {
	return r.waitResult().ErrorOrNil()
}

// Flat implements Result.
func (r *asyncOutput[T]) Flat(onOK func(T), onError func(error)) Output[T] {
	return r.waitResult().Flat(onOK, onError)
}

// GetOrDefault implements Result.
func (r *asyncOutput[T]) GetOrDefault(defaultValue T) T {
	return r.waitResult().GetOrDefault(defaultValue)
}

// GetOrNil implements Result.
func (r *asyncOutput[T]) GetOrNil() *T {
	return r.waitResult().GetOrNil()
}

// IfError implements Result.
func (r *asyncOutput[T]) IfError(onError func(error)) Output[T] {
	return r.waitResult().IfError(onError)
}

// IfOK implements Result.
func (r *asyncOutput[T]) IfOK(onOK func(T)) Output[T] {
	return r.waitResult().IfOK(onOK)
}

// IsError implements Result.
func (r *asyncOutput[T]) IsError() bool {
	return r.waitResult().IsError()
}

// IsOK implements Result.
func (r *asyncOutput[T]) IsOK() bool {
	return r.waitResult().IsOK()
}

func Async[T any](fn func() Output[T]) Output[T] {
	resultCh := make(chan Output[T])
	go func() {
		resultCh <- fn()
	}()
	return &asyncOutput[T]{resultCh: resultCh}
}

// Unwrap implements Result.
func (r *asyncOutput[T]) Unwrap() (T, error) {
	return r.waitResult().Unwrap()
}
