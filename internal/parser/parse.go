//go:generate go-wrap -mode=priv
package parser

import (
	"bytes"
	"go/ast"
	"go/printer"
	"os"
	"strconv"
	"strings"

	"github.com/catmorte/go-wrap/internal/declaration"

	. "github.com/catmorte/go-wrap/pkg/wrap"
	"golang.org/x/tools/go/packages"
)

type packageParser struct {
	*packages.Package
}

func loadPackages(path string, cfg *packages.Config) ([]*packages.Package, error) {
	return packages.Load(cfg, path+"/...")
}

func (p packageParser) extractRawCode(node any) string {
	b := new(bytes.Buffer)
	printer.Fprint(b, p.Fset, node)
	return b.String()
}

func (p packageParser) parseError(err packages.Error) *declaration.Error {
	return &declaration.Error{
		Message:  err.Msg,
		Position: err.Pos,
	}
}

func (p packageParser) newType(f *ast.Field) *declaration.Type[Empty] {
	return &declaration.Type[Empty]{Code: p.extractRawCode(f.Type)}
}

func (p packageParser) newTypeParams(f *ast.Field) *declaration.Type[declaration.ParamMeta] {
	code := p.extractRawCode(f.Type)
	return &declaration.Type[declaration.ParamMeta]{
		Code: code,
		Meta: declaration.ParamMeta{
			IsVararg: strings.HasPrefix(code, "..."),
		},
	}
}

func (p packageParser) newTypeTypeArg(f *ast.Field) *declaration.Type[declaration.TypeMeta] {
	return &declaration.Type[declaration.TypeMeta]{
		Code: p.extractRawCode(f.Type),
		Meta: declaration.TypeMeta{
			Name: f.Names[0].Name,
		},
	}
}

func newPackageParser(p *packages.Package) packageParser {
	return packageParser{p}
}

func unquote(v string) (string, error) {
	return strconv.Unquote(v)
}

func parseFields[T any](l *ast.FieldList, fn func(f *ast.Field) *declaration.Type[T]) []*declaration.Type[T] {
	if l == nil {
		return nil
	}
	res := make([]*declaration.Type[T], 0, len(l.List))
	for _, f := range l.List {
		res = append(res, fn(f))
	}
	return res
}

func filterFuncDecls(decls []ast.Decl) []*ast.FuncDecl {
	res := []*ast.FuncDecl{}
	for _, v := range decls {
		if f, ok := v.(*ast.FuncDecl); ok {
			res = append(res, f)
		}
	}
	return res
}

func (p packageParser) newFunc(fn *ast.FuncDecl) *declaration.Func {
	return &declaration.Func{
		Name:      fn.Name.Name,
		Code:      p.extractRawCode(fn),
		Params:    parseFields(fn.Type.Params, p.newTypeParams),
		Receivers: parseFields(fn.Recv, p.newType),
		Results:   parseFields(fn.Type.Results, p.newType),
		Types:     parseFields(fn.Type.TypeParams, p.newTypeTypeArg),
	}
}

func (p packageParser) newFile(fPath string, funcs []*declaration.Func, imports []*declaration.Import) *declaration.File {
	return &declaration.File{
		Path:    fPath,
		Funcs:   funcs,
		Imports: imports,
	}
}

func (p packageParser) newPackage(errors []*declaration.Error, files []*declaration.File) *declaration.Package {
	return &declaration.Package{
		Files:  files,
		Name:   p.Name,
		Errors: errors,
	}
}

func (p packageParser) newImport(v *ast.ImportSpec) (*declaration.Import, error) {
	pathUnquoted, err := unquote(v.Path.Value)
	if err != nil {
		return nil, err
	}
	alias := ""
	if v.Name != nil {
		alias = v.Name.String()
	}
	return &declaration.Import{Path: pathUnquoted, Alias: alias}, nil
}

func Parse(path string) []Out[*declaration.Package] {
	cfg := &packages.Config{
		Mode:  packages.NeedExportFile | packages.NeedModule | packages.NeedName | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedFiles,
		Dir:   ".",
		Env:   os.Environ(),
		Tests: false,
	}
	packagesLoaded := DisJoin(loadPackagesWrap(path, cfg))
	convertedToPackageParsers := EachAsync(packagesLoaded, newPackageParserWrap)
	return Each(convertedToPackageParsers, func(p packageParser) Out[*declaration.Package] {
		errorsParsed := EachAsync(OKVargs(p.Errors...), p.parseErrorWrap)
		filesDisJoined := DisJoin(OK(p.Syntax))
		filesParsed := EachAsync(filesDisJoined, func(f *ast.File) Out[*declaration.File] {
			fullPath := OK(p.Fset.Position(f.Pos()).Filename)
			funcsConverted := EachAsync(OKSlice(filterFuncDecls(f.Decls)), p.newFuncWrap)
			importsProcessed := EachAsync(OKSlice(f.Imports), p.newImportWrap)
			funcsJoined := JoinAsync(funcsConverted)
			importsJoined := JoinAsync(importsProcessed)
			return AndX3Async(fullPath, funcsJoined, importsJoined, p.newFileWrap)
		})
		errorsJoined := JoinAsync(errorsParsed)
		filesJoined := JoinAsync(filesParsed)
		return AndX2(errorsJoined, filesJoined, p.newPackageWrap)
	})
}
