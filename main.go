package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"

	. "github.com/catmorte/go-wrap/internal/declaration"
	"github.com/catmorte/go-wrap/internal/generator"
	"github.com/catmorte/go-wrap/internal/parser"
	. "github.com/catmorte/go-wrap/pkg/wrap"
)

var errEx = regexp.MustCompile(`"(.+)" imported and not used`)
var modes = map[string]func(*Func) bool{
	"pub": func(f *Func) bool {
		return unicode.IsUpper(rune(f.Name[0]))
	},
	"priv": func(f *Func) bool {
		return unicode.IsLower(rune(f.Name[0]))
	},
	"all": func(s *Func) bool {
		return true
	},
	"pub-rcv": func(f *Func) bool {
		return len(f.Receivers) > 0 && unicode.IsUpper(rune(f.Name[0]))
	},
	"priv-rcv": func(f *Func) bool {
		return len(f.Receivers) > 0 && unicode.IsLower(rune(f.Name[0]))
	},
	"all-rcv": func(f *Func) bool {
		return len(f.Receivers) > 0
	},
	"pub-fun": func(f *Func) bool {
		return len(f.Receivers) == 0 && unicode.IsUpper(rune(f.Name[0]))
	},
	"priv-fun": func(f *Func) bool {
		return len(f.Receivers) == 0 && unicode.IsLower(rune(f.Name[0]))
	},
	"all-fun": func(f *Func) bool {
		return len(f.Receivers) == 0
	},
}

func filterFuncs(fs []*Func, modeFunc func(*Func) bool, excludedFuncs []string) []*Func {
	res := []*Func{}
	for _, f := range fs {
		if modeFunc(f) && !slices.Contains(excludedFuncs, f.Name) && len(f.Results) <= 1 || (len(f.Results) == 2 && f.Results[1].Code == "error") {
			res = append(res, f)
		}
	}
	return res
}

func findFile(f string, vc []*File) *File {
	for _, v := range vc {
		if v.Path == f {
			return v
		}
	}
	return nil
}

func findImportByAlias(i []*Import, alias string) bool {
	for _, v := range i {
		if v.Alias == alias {
			return true
		}
	}
	return false
}

func findImportByPath(i []*Import, pkg string) bool {
	for _, v := range i {
		if v.Path == pkg && v.Alias != "" && v.Alias != "_" {
			return true
		}
	}
	return false
}

func findUnusedDotImportsForFile(filePath string, errors []*Error) []string {
	res := []string{}
	for _, e := range errors {
		if !strings.HasPrefix(e.Position, filePath) {
			continue
		}
		subMatches := errEx.FindStringSubmatch(e.Message)
		if len(subMatches) == 0 {
			continue
		}
		res = append(res, subMatches[1])
	}
	return res
}

func removeUnusedDotImports(unusedDotImports []string, imports []*Import) []*Import {
	res := []*Import{}
	for _, v := range imports {
		if v.Alias == "." && slices.Contains(unusedDotImports, v.Path) {
			continue
		}
		res = append(res, v)
	}
	return res
}

func findPackageAndFileByPath(fullPath string, ps []*Package) (*Package, *File) {
	var p *Package
	var f *File
	for _, v := range ps {
		f = findFile(fullPath, v.Files)
		if f == nil {
			continue
		}
		p = v
		break
	}
	return p, f
}

func main() {
	fileFlag := flag.String("file", "", "file")
	modeFlag := flag.String("mode", "all", "funcs/methods which needs to be wrapped: all/pub/priv[-rcv/-fun]")
	excludeFlag := flag.String("exclude", "", "coma separated list of funcs/methods to exclude")
	flag.Parse()
	file := os.Getenv("GOFILE")

	if file == "" {
		if fileFlag == nil || *fileFlag == "" {
			log.Fatal("file is not specified")
		}
		file = *fileFlag
	}

	modeFilter, ok := modes[*modeFlag]
	if !ok {
		log.Fatalf("unknown mode %v", modeFlag)
	}

	excludedFuncs := []string{}
	if excludeFlag != nil || *excludeFlag != "" {
		excludedFuncs = strings.Split(*excludeFlag, ",")
	}

	pathGot := Wrap(os.Getwd())
	fileSaved := And(pathGot, func(path string) Out[string] {
		fullPath := filepath.Join(path, file)
		packagesParsed := parser.Parse(path)
		packagesJoined := JoinAsync(packagesParsed)
		return And(packagesJoined, func(ps []*Package) Out[string] {
			p, f := findPackageAndFileByPath(fullPath, ps)
			if p == nil || f == nil {
				return Err[string](fmt.Errorf("file %v not found", fullPath))
			}
			f.Funcs = filterFuncs(f.Funcs, modeFilter, excludedFuncs)
			ok := findImportByPath(f.Imports, WrapPkgPath)
			if !ok {
				i := 0
				alias := ""
				for {
					alias = fmt.Sprintf("%s%d", WrapPkgAlias, i)
					ok = findImportByAlias(f.Imports, alias)
					if !ok {
						f.Imports = append(f.Imports, &Import{Alias: alias, Path: WrapPkgPath})
						break
					}
					i++
				}
			}
			fileName := fmt.Sprintf("%s_wrap.go", strings.TrimSuffix(file, filepath.Ext(file)))
			fullPath := filepath.Join(path, fileName)
			codeGenerated := generator.Generate(p.Name, *f, false)
			return And(codeGenerated, func(raw []byte) Out[string] {
				return Wrap(fullPath, os.WriteFile(fullPath, raw, 0644))
			})
		})
	})
	AndX2(pathGot, fileSaved, func(path, filePath string) Out[Empty] {
		packagesParsed := parser.Parse(path)
		packagesJoined := JoinAsync(packagesParsed)
		return And(packagesJoined, func(ps []*Package) Out[Empty] {
			p, f := findPackageAndFileByPath(filePath, ps)
			if p == nil || f == nil {
				return Err[Empty](fmt.Errorf("file %v not found", filePath))
			}
			unusedImports := findUnusedDotImportsForFile(filePath, p.Errors)
			f.Imports = removeUnusedDotImports(unusedImports, f.Imports)

			codeGenerated := generator.Generate(p.Name, *f, true)
			return And(codeGenerated, func(raw []byte) Out[Empty] {
				return Void(os.WriteFile(filePath, raw, 0644))
			})
		})
	}).IfError(func(err error) {
		log.Fatal(err)
	})
}
