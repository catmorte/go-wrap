package wrap

import (
	"errors"
)

var (
	ErrChanClosed = errors.New("channel closed")
	ErrNotFound   = errors.New("unable to find first, condition not met")
)

type (
	Successor[T1, TT any]                                   func(T1) Out[TT]
	SuccessorX2[T1, T2, TT any]                             func(T1, T2) Out[TT]
	SuccessorX3[T1, T2, T3, TT any]                         func(T1, T2, T3) Out[TT]
	SuccessorX4[T1, T2, T3, T4, TT any]                     func(T1, T2, T3, T4) Out[TT]
	SuccessorX5[T1, T2, T3, T4, T5, TT any]                 func(T1, T2, T3, T4, T5) Out[TT]
	SuccessorX6[T1, T2, T3, T4, T5, T6, TT any]             func(T1, T2, T3, T4, T5, T6) Out[TT]
	SuccessorX7[T1, T2, T3, T4, T5, T6, T7, TT any]         func(T1, T2, T3, T4, T5, T6, T7) Out[TT]
	SuccessorX8[T1, T2, T3, T4, T5, T6, T7, T8, TT any]     func(T1, T2, T3, T4, T5, T6, T7, T8) Out[TT]
	SuccessorX9[T1, T2, T3, T4, T5, T6, T7, T8, T9, TT any] func(T1, T2, T3, T4, T5, T6, T7, T8, T9) Out[TT]
)

func Proof(r ...ErrorContainer) Out[Empty] {
	for _, v := range r {
		if err := v.ErrorOrNil(); err != nil {
			return Err[Empty](err)
		}
	}
	return OK(Empty{})
}

func ProofAsync(r ...ErrorContainer) Out[Empty] {
	return Async(func() Out[Empty] {
		return Proof(r...)
	})
}

func first[T any](r []Out[T], test func(Out[T]) bool) Out[T] {
	length := len(r)
	chans := make(chan Out[T], length)
	for _, v := range r {
		v := v
		go func() {
			_ = v.IsOK()
			chans <- v
		}()
	}
	for v := range chans {
		if test(v) {
			return v
		}
		if length == 0 {
			break
		}
	}
	return Err[T](ErrNotFound)
}

func FirstOK[T any](r []Out[T]) Out[T] {
	return first(r, func(o Out[T]) bool {
		return o.IsOK()
	})
}

func FirstErr[T any](r []Out[T]) Out[T] {
	return first(r, func(o Out[T]) bool {
		return o.IsError()
	})
}

func OnlyOKs[T any](r []Out[T]) []Out[T] {
	res := make([]Out[T], 0, len(r))
	for _, v := range r {
		if v.IsOK() {
			res = append(res, v)
		}
	}
	return res
}

func OnlyErrors[T any](r []Out[T]) []Out[T] {
	res := make([]Out[T], 0, len(r))
	for _, v := range r {
		if v.IsError() {
			res = append(res, v)
		}
	}
	return res
}

func eachFunc[T any, TT any](r []Out[T], f Successor[T, TT], and func(Out[T], Successor[T, TT]) Out[TT]) []Out[TT] {
	res := make([]Out[TT], 0, len(r))
	for _, v := range r {
		res = append(res, and(v, f))
	}
	return res
}

func Each[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	return eachFunc(r, f, And)
}

func EachAsync[T any, TT any](r []Out[T], f Successor[T, TT]) []Out[TT] {
	return eachFunc(r, f, AndAsync)
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
	res := make([]Out[TT], 0, n)
	for i := 0; i < n; i++ {
		res = append(res, and(OK(i), f))
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
	res := make([]T, 0, len(r))
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
		res := make([]Out[T], 0, len(values))
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
	return OK(val)
}

func Void(err error) Out[Empty] {
	return Wrap(Empty{}, err)
}

func VoidVargs(values ...error) []Out[Empty] {
	res := make([]Out[Empty], 0, len(values))
	for _, v := range values {
		res = append(res, Void(v))
	}
	return res
}

func Flat[T any](r []Out[[]T]) []Out[T] {
	res := []Out[T]{}
	for _, v := range r {
		res = append(res, DisJoin(v)...)
	}
	return res
}

func And[T any, TT any](r Out[T], f Successor[T, TT]) Out[TT] {
	if r.IsOK() {
		var defaultV T
		return f(r.GetOrDefault(defaultV))
	}
	return Err[TT](r.ErrorOrNil())
}

func AndAsync[T any, TT any](r Out[T], f Successor[T, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return And(r, f)
	})
}

func AndX2[T1, T2, TT any](r1 Out[T1], r2 Out[T2], f SuccessorX2[T1, T2, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
		)
	})
}

func AndX2Async[T1, T2, TT any](r1 Out[T1], r2 Out[T2], f SuccessorX2[T1, T2, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX2(r1, r2, f)
	})
}

func AndX3[T1, T2, T3, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], f SuccessorX3[T1, T2, T3, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
		)
	})
}

func AndX3Async[T1, T2, T3, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], f SuccessorX3[T1, T2, T3, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX3(r1, r2, r3, f)
	})
}

func AndX4[T1, T2, T3, T4, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], f SuccessorX4[T1, T2, T3, T4, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		var defaultV4 T4
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
			r4.GetOrDefault(defaultV4),
		)
	})
}

func AndX4Async[T1, T2, T3, T4, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], f SuccessorX4[T1, T2, T3, T4, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX4(r1, r2, r3, r4, f)
	})
}

func AndX5[T1, T2, T3, T4, T5, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], f SuccessorX5[T1, T2, T3, T4, T5, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		var defaultV4 T4
		var defaultV5 T5
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
			r4.GetOrDefault(defaultV4),
			r5.GetOrDefault(defaultV5),
		)
	})
}

func AndX5Async[T1, T2, T3, T4, T5, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], f SuccessorX5[T1, T2, T3, T4, T5, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX5(r1, r2, r3, r4, r5, f)
	})
}

func AndX6[T1, T2, T3, T4, T5, T6, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], f SuccessorX6[T1, T2, T3, T4, T5, T6, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		var defaultV4 T4
		var defaultV5 T5
		var defaultV6 T6
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
			r4.GetOrDefault(defaultV4),
			r5.GetOrDefault(defaultV5),
			r6.GetOrDefault(defaultV6),
		)
	})
}

func AndX6Async[T1, T2, T3, T4, T5, T6, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], f SuccessorX6[T1, T2, T3, T4, T5, T6, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX6(r1, r2, r3, r4, r5, r6, f)
	})
}

func AndX7[T1, T2, T3, T4, T5, T6, T7, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], r7 Out[T7], f SuccessorX7[T1, T2, T3, T4, T5, T6, T7, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		var defaultV4 T4
		var defaultV5 T5
		var defaultV6 T6
		var defaultV7 T7
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
			r4.GetOrDefault(defaultV4),
			r5.GetOrDefault(defaultV5),
			r6.GetOrDefault(defaultV6),
			r7.GetOrDefault(defaultV7),
		)
	})
}

func AndX7Async[T1, T2, T3, T4, T5, T6, T7, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], r7 Out[T7], f SuccessorX7[T1, T2, T3, T4, T5, T6, T7, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX7(r1, r2, r3, r4, r5, r6, r7, f)
	})
}

func AndX8[T1, T2, T3, T4, T5, T6, T7, T8, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], r7 Out[T7], r8 Out[T8], f SuccessorX8[T1, T2, T3, T4, T5, T6, T7, T8, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		var defaultV4 T4
		var defaultV5 T5
		var defaultV6 T6
		var defaultV7 T7
		var defaultV8 T8
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
			r4.GetOrDefault(defaultV4),
			r5.GetOrDefault(defaultV5),
			r6.GetOrDefault(defaultV6),
			r7.GetOrDefault(defaultV7),
			r8.GetOrDefault(defaultV8),
		)
	})
}

func AndX8Async[T1, T2, T3, T4, T5, T6, T7, T8, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], r7 Out[T7], r8 Out[T8], f SuccessorX8[T1, T2, T3, T4, T5, T6, T7, T8, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX8(r1, r2, r3, r4, r5, r6, r7, r8, f)
	})
}

func AndX9[T1, T2, T3, T4, T5, T6, T7, T8, T9, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], r7 Out[T7], r8 Out[T8], r9 Out[T9], f SuccessorX9[T1, T2, T3, T4, T5, T6, T7, T8, T9, TT]) Out[TT] {
	allGood := Proof(r1, r2)
	return And(allGood, func(Empty) Out[TT] {
		var defaultV1 T1
		var defaultV2 T2
		var defaultV3 T3
		var defaultV4 T4
		var defaultV5 T5
		var defaultV6 T6
		var defaultV7 T7
		var defaultV8 T8
		var defaultV9 T9
		return f(r1.GetOrDefault(defaultV1),
			r2.GetOrDefault(defaultV2),
			r3.GetOrDefault(defaultV3),
			r4.GetOrDefault(defaultV4),
			r5.GetOrDefault(defaultV5),
			r6.GetOrDefault(defaultV6),
			r7.GetOrDefault(defaultV7),
			r8.GetOrDefault(defaultV8),
			r9.GetOrDefault(defaultV9),
		)
	})
}

func AndX9Async[T1, T2, T3, T4, T5, T6, T7, T8, T9, TT any](r1 Out[T1], r2 Out[T2], r3 Out[T3], r4 Out[T4], r5 Out[T5], r6 Out[T6], r7 Out[T7], r8 Out[T8], r9 Out[T9], f SuccessorX9[T1, T2, T3, T4, T5, T6, T7, T8, T9, TT]) Out[TT] {
	return Async(func() Out[TT] {
		return AndX9(r1, r2, r3, r4, r5, r6, r7, r8, r9, f)
	})
}
