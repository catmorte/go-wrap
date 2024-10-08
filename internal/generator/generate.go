//go:generate go-wrap -mode=priv
package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/catmorte/go-wrap/internal/declaration"
	. "github.com/catmorte/go-wrap/pkg/wrap"
	"golang.org/x/tools/imports"
)

const fileTemplate = `
// Code generated by "go-wrap"; DO NOT EDIT.
package {{.PackageName}}
import (
{{range .Imports -}}
		{{if .Alias }} {{.Alias }} {{end -}}
		"{{.Path}}"
{{end -}}
)
{{if .PrintRaw}}
{{range .Funcs}}
{{.Code}}
{{end}}
{{ else}}
{{range .Funcs}}
func {{if .Receivers}}({{range .Receivers}}rcv {{.Code}}{{end}}){{end -}}
{{.Name}}Wrap{{- if .Types}}[{{range $index, $param := .Types}}{{.Meta.Name}} {{.Code}} {{if $index}}, {{end}} {{end}}]{{end -}}
({{range $index, $param := .Params}}{{if $index}}, {{end}}arg{{$index}} {{.Code}}{{end}}) {{"" -}} 
{{if eq (len .Results) 2 -}}
	{{$.WrapPackageAlias}}Out[{{(index .Results 0).Code}}] {
	return {{$.WrapPackageAlias}}Wrap({{if .Receivers}}{{range .Receivers}}rcv.{{end}}{{end}}{{.Name}}({{range $index, $param := .Params}}{{if $index}}, {{end}}arg{{$index}}{{if $param.Meta.IsVararg}}...{{end}} {{end}}))
{{- else if eq (len .Results) 1 -}}
	{{if eq (index .Results 0).Code "error" -}}
		{{$.WrapPackageAlias}}Out[{{$.WrapPackageAlias}}Empty] {
		return {{$.WrapPackageAlias}}Void({{if .Receivers}}{{range .Receivers}}rcv.{{end}}{{end}}{{.Name}}({{range $index, $param := .Params}}{{if $index}}, {{end}}arg{{$index}}{{if $param.Meta.IsVararg}}...{{end}} {{end}}))
	{{- else -}}
		{{$.WrapPackageAlias}}Out[{{(index .Results 0).Code}}] {
		return {{$.WrapPackageAlias}}OK({{if .Receivers}}{{range .Receivers}}rcv.{{end}}{{end}}{{.Name}}({{range $index, $param := .Params}}{{if $index}}, {{end}}arg{{$index}}{{if $param.Meta.IsVararg}}...{{end}} {{end}}))
	{{- end }}
{{- else if eq (len .Results) 0 -}}
	{{$.WrapPackageAlias}}Out[{{$.WrapPackageAlias}}Empty] {
	{{if .Receivers}}{{range .Receivers}}rcv.{{end}}{{end}}{{.Name}}({{range $index, $param := .Params}}{{if $index}}, {{end}}arg{{$index}}{{if $param.Meta.IsVararg}}...{{end}} {{end}})
	return {{$.WrapPackageAlias}}OK({{$.WrapPackageAlias}}Empty{})
{{- end }} 
}
{{end}}
{{end}}
`

type fileTemplateData struct {
	PackageName      string
	WrapPackageAlias string
	PrintRaw         bool
	declaration.File
}

func getWrapPrefix(imports []*declaration.Import) string {
	var imp *declaration.Import
	for _, i := range imports {
		if i.Path == declaration.WrapPkgPath && i.Alias != "_" && i.Alias != "" {
			imp = i
		}
	}
	if imp == nil || imp.Alias == "." {
		return ""
	}
	return fmt.Sprintf("%s.", imp.Alias)
}

func parseTemplate() (*template.Template, error) {
	return template.New("").Parse(fileTemplate)
}

func executeTemplate[T any](t *template.Template, data T) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func newFileTemplateData(prefix string, packageName string, f declaration.File, printRaw bool) fileTemplateData {
	return fileTemplateData{
		PackageName:      packageName,
		File:             f,
		WrapPackageAlias: prefix,
		PrintRaw:         printRaw,
	}
}

func formatSource(b []byte) ([]byte, error) {
	return format.Source(b)
}

func formatImports(b []byte) ([]byte, error) {
	return imports.Process("", b, nil)
}

func Generate(packageName string, f declaration.File, printRaw bool) Out[[]byte] {
	gotWrapPreifx := getWrapPrefixWrap(f.Imports)
	templateDataCreated := AndX4Async(gotWrapPreifx, OK(packageName), OK(f), OK(printRaw), newFileTemplateDataWrap)
	templateParsed := parseTemplateWrap()
	return AndX2Async(templateParsed, templateDataCreated, func(t *template.Template, data fileTemplateData) Out[[]byte] {
		codeGenerated := executeTemplateWrap(t, data)
		codeFormatted := AndAsync(codeGenerated, formatSourceWrap)
		return AndAsync(codeFormatted, formatImportsWrap)
	})
}
