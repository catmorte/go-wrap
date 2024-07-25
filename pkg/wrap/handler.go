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

func Proof(r ...ErrorContainer) OutVoid {
	for _, v := range r {
		if err := v.ErrorOrNil(); err != nil {
			return Err[Void](err)
		}
	}
	return OK(Void{})
}

func ProofAsync(r ...ErrorContainer) OutVoid {
	return Async[Void](func() Out[Void] {
		return Proof(r...)
	})
}

func ForEach[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	res := make([]Out[TT], len(r))
	for _, v := range r {
		res = append(res, And(v, f))
	}
	return res
}

func ForEachAsync[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	res := make([]Out[TT], len(r))
	for _, v := range r {
		res = append(res, AndAsync(v, f))
	}
	return res
}

func ForRange[TT any](n int, f Successor[int, TT]) []Out[TT] {
	res := make([]Out[TT], n)
	for i := 0; i < n; i++ {
		res = append(res, f(i))
	}
	return res
}

func ForRangeAsync[TT any](n int, f Successor[int, TT]) []Out[TT] {
	res := make([]Out[TT], n)
	for i := 0; i < n; i++ {
		res = append(res, AndAsync(OK(i), f))
	}
	return res
}

func Join[T any, TT any](r []Out[T], f SuccessorSlice[T, TT]) Out[[]TT] {
	res := make([]T, len(r))
	for _, v := range r {
		if v.IsError() {
			return Err[[]TT](v.ErrorOrNil())
		}
		var defaultV T
		res = append(res, v.GetOrDefault(defaultV))
	}
	return f(res)
}

func JoinAsync[T any, TT any](r []Out[T], f SuccessorSlice[T, TT]) Out[[]TT] {
	return Async[[]TT](func() Out[[]TT] {
		return Join(r, f)
	})
}

func Unjoin[T any, TT any](r Out[[]T], f Successor[T, TT]) []Out[TT] {
	if r.IsOK() {
		var defaultV []T
		values := r.GetOrDefault(defaultV)
		res := make([]Out[TT], len(values))
		for _, v := range values {
			res = append(res, f(v))
		}
		return res
	}
	return []Out[TT]{Err[TT](r.ErrorOrNil())}
}

func UnjoinAsync[T any, TT any](r Out[[]T], f Successor[T, TT]) []Out[TT] {
	if r.IsOK() {
		var defaultV []T
		values := r.GetOrDefault(defaultV)
		res := make([]Out[TT], len(values))
		for _, v := range values {
			res = append(res, Async(func() Out[TT] {
				return f(v)
			}))
		}
		return res
	}
	return []Out[TT]{Err[TT](r.ErrorOrNil())}
}

func Wrap[T any](val T, err error) Out[T] {
	if err != nil {
		return Err[T](err)
	}
	return OK[T](val)
}

func WrapVoid(err error) OutVoid {
	return Wrap(Void{}, err)
}
