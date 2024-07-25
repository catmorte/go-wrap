package wrap

type (
	Successor[T any, TT any]      func(T) Out[TT]
	SuccessorSlice[T any, TT any] func([]T) Out[[]TT]
)

func And[T any, TT any](r Out[T], f Successor[T, TT]) Out[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}

func AndAsync[T any, TT any](r Out[T], f Successor[T, TT]) Out[TT] {
	return Async[TT](func() Out[TT] {
		return And[T, TT](r, f)
	})
}

func Proof(r ...ErrorContainer) Out[Empty] {
	for _, v := range r {
		if err := v.ErrorOrNil(); err != nil {
			return Err[Empty](err)
		}
	}
	return OK(Empty{})
}

func ProofAsync(r ...ErrorContainer) Out[Empty] {
	return Async[Empty](func() Out[Empty] {
		return Proof(r...)
	})
}

func Each[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	res := make([]Out[TT], len(r))
	for _, v := range r {
		res = append(res, And(v, f))
	}
	return res
}

func EachAsync[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	res := make([]Out[TT], len(r))
	for _, v := range r {
		res = append(res, AndAsync(v, f))
	}
	return res
}

func Range[TT any](n int, f Successor[int, TT]) []Out[TT] {
	res := make([]Out[TT], n)
	for i := 0; i < n; i++ {
		res = append(res, f(i))
	}
	return res
}

func RangeAsync[TT any](n int, f Successor[int, TT]) []Out[TT] {
	res := make([]Out[TT], n)
	for i := 0; i < n; i++ {
		res = append(res, AndAsync(OK(i), f))
	}
	return res
}

func Join[T any](r []Out[T]) Out[[]T] {
	res := make([]T, len(r))
	for _, v := range r {
		if v.IsError() {
			return Err[[]T](v.ErrorOrNil())
		}
		var defaultV T
		res = append(res, v.GetOrDefault(defaultV))
	}
	return OK(res)
}

func DisJoin[T any](r Out[[]T]) []Out[T] {
	if r.IsOK() {
		var defaultV []T
		values := r.GetOrDefault(defaultV)
		res := make([]Out[T], len(values))
		for _, v := range values {
			res = append(res, OK(v))
		}
		return res
	}
	return []Out[T]{Err[T](r.ErrorOrNil())}
}

func Wrap[T any](val T, err error) Out[T] {
	if err != nil {
		return Err[T](err)
	}
	return OK[T](val)
}

func Void(err error) Out[Empty] {
	return Wrap(Empty{}, err)
}
