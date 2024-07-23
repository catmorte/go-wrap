package wrap

type (
	Successor[T any, TT any]      func(T) Output[TT]
	SuccessorSlice[T any, TT any] func([]T) Output[[]TT]
)

func Next[T any, TT any](f Successor[T, TT], r Output[T]) Output[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}

func NextAsync[T any, TT any](f Successor[T, TT], r Output[T]) Output[TT] {
	return Async[TT](func() Output[TT] {
		return Next[T, TT](f, r)
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

func ForEach[T any, TT any](f Successor[T, TT], r ...Output[T]) []Output[TT] {
	res := make([]Output[TT], len(r))
	for _, v := range r {
		res = append(res, Next(f, v))
	}
	return res
}

func ForEachAsync[T any, TT any](f Successor[T, TT], r ...Output[T]) []Output[TT] {
	res := make([]Output[TT], len(r))
	for _, v := range r {
		res = append(res, NextAsync(f, v))
	}
	return res
}

func Join[T any, TT any](f SuccessorSlice[T, TT], r ...Output[T]) Output[[]TT] {
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

func JoinAsync[T any, TT any](f SuccessorSlice[T, TT], r ...Output[T]) Output[[]TT] {
	return Async[[]TT](func() Output[[]TT] {
		return Join(f, r...)
	})
}

func DisJoin[T any, TT any](f Successor[T, TT], r Output[[]T]) []Output[TT] {
	if r.IsOK() {
		var defaultV []T
		values := r.GetOrDefault(defaultV)
		res := make([]Output[TT], len(values))
		for _, v := range values {
			res = append(res, f(v))
		}
		return res
	}
	return []Output[TT]{Err[TT](r.ErrorOrNil())}
}

func DisJoinAsync[T any, TT any](f Successor[T, TT], r Output[[]T]) []Output[TT] {
	if r.IsOK() {
		var defaultV []T
		values := r.GetOrDefault(defaultV)
		res := make([]Output[TT], len(values))
		for _, v := range values {
			res = append(res, Async(func() Output[TT] {
				return f(v)
			}))
		}
		return res
	}
	return []Output[TT]{Err[TT](r.ErrorOrNil())}
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
