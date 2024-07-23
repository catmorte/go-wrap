# go-wrap

## probably it's just some sort of brainfuck, but still pretty funny concept imho.

So, the concept is to use `Output[T]` for each func within the project, or to use `Wrap[T](val T, err error)` (yet another function in this package to wrap common touples `return val, err`), or just `OK[T](v T)` / `Err[T](err error)` to convert any value to it. 
togeather with the list of handlers such as `Continue`, `Proof`, `ContinueForEach`, `ContinueSlice` (plus `**Async` versions of those functions) to reach some kind of declarative style.
