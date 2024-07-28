package wrap

import "errors"

var ErrChanClosed = errors.New("channel closed")

type (
	Successor[T any, TT any] func(T) Out[TT]
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

func eachFunc[T any, TT any](r []Out[T], f Successor[T, TT], and func(Out[T], Successor[T, TT]) Out[TT]) []Out[TT] {
	res := make([]Out[TT], len(r))
	for _, v := range r {
		res = append(res, and(v, f))
	}
	return res
}

func Each[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	return eachFunc[T, TT](r, f, And)
}

func EachAsync[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	return eachFunc[T, TT](r, f, AndAsync)
}

func slicedFunc[T any, TT any](
	sliceSize int,
	r []Out[T],
	f Successor[[]T, []TT],
	and func(Out[[]T], Successor[[]T, []TT]) Out[[]TT],
	join func(r []Out[T]) Out[[]T],
) []Out[[]TT] {
	res := []Out[[]TT]{}
	for i := 0; i < len(r); i += sliceSize {
		end := i + sliceSize
		if end > len(r) {
			end = len(r)
		}
		res = append(res, and(join(r[i:end]), f))
	}
	return res
}

func Sliced[T any, TT any](sliceSize int, r []Out[T], f Successor[[]T, []TT]) []Out[[]TT] {
	return slicedFunc(sliceSize, r, f, And, Join)
}

func SlicedAsync[T any, TT any](sliceSize int, r []Out[T], f Successor[[]T, []TT]) []Out[[]TT] {
	return slicedFunc(sliceSize, r, f, AndAsync, JoinAsync)
}

func rangeFunc[TT any](n int, f Successor[int, TT], and func(Out[int], Successor[int, TT]) Out[TT]) []Out[TT] {
	res := make([]Out[TT], n)
	for i := 0; i < n; i++ {
		res = append(res, And(OK(i), f))
	}
	return res
}

func Range[TT any](n int, f Successor[int, TT]) []Out[TT] {
	return rangeFunc(n, f, And)
}

func RangeAsync[TT any](n int, f Successor[int, TT]) []Out[TT] {
	return rangeFunc(n, f, AndAsync)
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

func JoinAsync[T any](r []Out[T]) Out[[]T] {
	return Async[[]T](func() Out[[]T] {
		return Join(r)
	})
}

func Just[T any](isError bool, r []Out[T]) []Out[T] {
	res := []Out[T]{}
	for _, v := range r {
		if v.IsError() == isError {
			continue
		}
		res = append(res, v)
	}
	return res
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

func ReadChan[T any](ch <-chan T, onClosed func() Out[T]) Out[T] {
	v, ok := <-ch
	if ok {
		return OK(v)
	}
	return onClosed()
}

func ReadChanAsync[T any](ch <-chan T, onClosed func() Out[T]) Out[T] {
	return Async[T](func() Out[T] {
		return ReadChan(ch, onClosed)
	})
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

func VoidVargs(values ...error) []Out[Empty] {
	res := make([]Out[Empty], len(values))
	for _, v := range values {
		res = append(res, Void(v))
	}
	return res
}
