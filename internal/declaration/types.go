package declaration

import . "github.com/catmorte/go-wrap/pkg/wrap"

type (
	Import struct {
		Alias string
		Path  string
	}
	Type[T any] struct {
		Code string
		Meta T
	}
	TypeMeta struct {
		Name string
	}
	ParamMeta struct {
		IsVararg bool
	}
	Func struct {
		Name      string
		Code      string
		Params    []*Type[ParamMeta]
		Results   []*Type[Empty]
		Types     []*Type[TypeMeta]
		Receivers []*Type[Empty]
	}
	File struct {
		Path    string
		Imports []*Import
		Funcs   []*Func
	}
	Package struct {
		Name   string
		Files  []*File
		Errors []*Error
	}
	Error struct {
		Message  string
		Position string
	}
)
