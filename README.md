# go-wrap

## probably it's just some sort of brainfuck, but still pretty funny concept imho.

So, the concept is to use `Out[T]` for each func within the project, or to use `Wrap[T](val T, err error)` (yet another function in this package to wrap common touples `return val, err`), or just `OK[T](v T)` / `Err[T](err error)` to convert any value to it.

Togeather with the list of handlers such as `And` (also `AndXN` where N is a number up to 9), `Join`, `Proof`, `Range`, `Each`, `ReadChan`, `Sliced` (plus `**Async` versions of those functions) and also `Just`,`DisJoin`, `Flat` to reach some kind of declarative style.

## go generate

U can `go install` this package and later use it via `//go:generate go-wrap` togeather with one of the flags `public`(default), `private` and `all` to generate wrappers for regular functions which can return few values like:

- ()
- (V)
- (error)
- (V, error)

where V is value of any type
