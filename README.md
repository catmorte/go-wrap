# go-wrap

## Overview

The concept of this package is to use `Out[T]` for each function within the project or to use `Wrap[T](val T, err error)` (another function in this package to wrap common tuples like `return val, err`). Alternatively, you can use `OK[T](v T)` or `Err[T](err error)` to convert any value into these forms.

The package includes a list of handlers such as `And` (also `AndXN` where N is a number up to 9), `Join`, `Proof`, `Range`, `Each`, `ReadChan`, `Sliced` (with `Async` versions of these functions). Additionally, `Just`, `DisJoin`, and `Flat` are provided to create a more declarative style.

## go generate

You can `go install` this package and later use it with `//go:generate go-wrap`, along with the following flags:

- **exclude**: A list of comma-separated function names to exclude.
- **mode**: Specifies the visibility and type of functions to generate. Available options are:
  - `pub` (default)
  - `priv`
  - `all`
  - `priv-rcv`
  - `pub-rcv`
  - `all-rcv`
  - `priv-fun`
  - `pub-fun`
  - `all-fun`

These flags are used to generate wrappers for regular functions and methods that can return a variety of signatures:

- ()
- (V)
- (error)
- (V, error)

Where `V` is a value of any type.
