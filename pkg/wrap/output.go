package wrap

type ErrorContainer interface {
	ErrorOrNil() error
	IsOK() bool
	IsError() bool
}

type Void struct{}

type Output[T any] interface {
	ErrorContainer

	GetOrDefault(defaultValue T) T
	GetOrNil() *T

	IfOK(onOk func(T)) Output[T]
	IfError(onError func(error)) Output[T]
	Flat(onOK func(T), onError func(error)) Output[T]

	Unwrap() (T, error)
}
