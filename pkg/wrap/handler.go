package wrap

type (
	Successor[T any, TT any]      func(T) Output[TT]
	SuccessorSlice[T any, TT any] func([]T) Output[[]TT]
)

func Continue[T any, TT any](f Successor[T, TT], r Output[T]) Output[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}

func ContinueAsync[T any, TT any](f Successor[T, TT], r Output[T]) Output[TT] {
	return Async[TT](func() Output[TT] {
		return Continue[T, TT](f, r)
	})
}

func Proof(r ...ErrorContainer) Output[Void] {
	for _, v := range r {
		if err := v.ErrorOrNil(); err != nil {
			return Err[Void](err)
		}
	}
	return OK(Void{})
}

func ProofAsync(r ...ErrorContainer) Output[Void] {
	return Async[Void](func() Output[Void] {
		return Proof(r...)
	})
}

func ContinueForEach[T any, TT any](f Successor[T, TT], r ...Output[T]) []Output[TT] {
	var res []Output[TT]
	for _, v := range r {
		res = append(res, Continue(f, v))
	}
	return res
}

func ContinueForEachAsync[T any, TT any](f Successor[T, TT], r ...Output[T]) []Output[TT] {
	var res []Output[TT]
	for _, v := range r {
		res = append(res, ContinueAsync(f, v))
	}
	return res
}

func ContinueSlice[T any, TT any](f SuccessorSlice[T, TT], r ...Output[T]) Output[[]TT] {
	var res []T
	for _, v := range r {
		if v.IsError() {
			return Err[[]TT](v.ErrorOrNil())
		}
		var defaultV T
		res = append(res, v.GetOrDefault(defaultV))
	}
	return f(res)
}

func ContinueSliceAsync[T any, TT any](f SuccessorSlice[T, TT], r ...Output[T]) Output[[]TT] {
	return Async[[]TT](func() Output[[]TT] {
		return ContinueSlice(f, r...)
	})
}

func Wrap[T any](val T, err error) Output[T] {
	if err != nil {
		return Err[T](err)
	}
	return OK[T](val)
}

func WrapVoid(err error) Output[Void] {
	if err != nil {
		return Err[Void](err)
	}
	return OK[Void](Void{})
}
