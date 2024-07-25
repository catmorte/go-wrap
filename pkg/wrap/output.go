package wrap

type (
	Empty struct{}

	ErrorContainer interface {
		ErrorOrNil() error
		IsOK() bool
		IsError() bool
	}

	Out[T any] interface {
		ErrorContainer
		GetOrDefault(defaultValue T) T
		GetOrNil() *T
		IfOK(onOk func(T)) Out[T]
		IfError(onError func(error)) Out[T]
		Flat(onOK func(T), onError func(error)) Out[T]
		Unwrap() (T, error)
	}
)
