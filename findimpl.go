package findimpl

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const doc = "findimpl is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "findimpl",
	Doc:  doc,
	Run:  run,
}

var target string // -target=io.Writer

func init() {
	Analyzer.Flags.StringVar(&target, "target", target, "")
}

func run(pass *analysis.Pass) (interface{}, error) {
	typeSpecs := []*ast.TypeSpec{}

	for _, f := range pass.Files {
		for _, d := range f.Decls {
			genDecl, _ := d.(*ast.GenDecl)
			if genDecl == nil {
				continue
			}
			for _, s := range genDecl.Specs {
				typeSpec, _ := s.(*ast.TypeSpec)
				if typeSpec == nil {
					continue
				}
				typeSpecs = append(typeSpecs, typeSpec)
			}
		}
	}

	return runTypeSpecs(pass, typeSpecs)
}

func runTypeSpecs(pass *analysis.Pass, typeSpecs []*ast.TypeSpec) (interface{}, error) {
	targetInterface, err := getInterface(pass)
	if err != nil {
		return nil, err
	}

	for _, ts := range typeSpecs {
		if implements(pass.TypesInfo.TypeOf(ts.Name), targetInterface) {
			pass.Reportf(ts.Name.Pos(), "%s implements %s", ts.Name.Name, target)
		}
	}
	return nil, nil
}

func implements(V types.Type, T *types.Interface) bool {
	if types.Implements(V, T) {
		return true
	}

	if types.Implements(types.NewPointer(V), T) {
		return true
	}

	return false
}

// getInterface gets *types.Interface by target
func getInterface(pass *analysis.Pass) (*types.Interface, error) {
	if target == "error" {
		return types.Universe.Lookup("error").Type().Underlying().(*types.Interface), nil
	}

	if !strings.Contains(target, ".") {
		return nil, fmt.Errorf("invalid target: %s", target)
	}

	// target -> targetImportPath, targetInterfaceName
	// "io.Writer" -> "io", "Writer"
	// `"hoge.fuga/piyo".Foo` -> `"hoge.fuga/piyo"`, `Foo`
	lastComma := strings.LastIndex(target, ".")
	targetImportPath := target[:lastComma]
	targetInterfaceName := target[lastComma+1:]

	buildPkg, err := build.Default.Import(targetImportPath, ".", build.ImportMode(0))
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, buildPkg.Dir, nil, parser.Mode(0))
	if err != nil {
		return nil, err
	}
	// hoge
	if _, ok := pkgs[buildPkg.Name]; !ok {
		return nil, errors.New("unexpected")
	}
	pkg := pkgs[buildPkg.Name]

	files := make([]*ast.File, 0, len(pkg.Files))
	for _, f := range pkg.Files {
		files = append(files, f)
	}

	// 型情報を持ってくる
	c := &types.Config{
		Importer: importer.Default(),
	}
	info := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{},
	}
	if _, err := c.Check(buildPkg.ImportPath, fset, files, info); err != nil {
		return nil, err
	}

	for _, tv := range info.Types {
		named, _ := tv.Type.(*types.Named)
		if named == nil {
			continue
		}
		if named.Obj().Name() != targetInterfaceName {
			continue
		}

		i, _ := named.Underlying().(*types.Interface)
		if i == nil {
			return nil, fmt.Errorf("%s is not interface", target)
		}
		return i, nil
	}

	return nil, fmt.Errorf("%s not found", target)
}
