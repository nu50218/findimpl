package findimpl

import (
	"go/ast"
	"go/types"

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
	targetInterface, err := getInterface()
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
func getInterface() (*types.Interface, error) {
	// とりあえずエラーを満たすやつを見つけるコードがかけた
	// TODO: targetに応じて返す*types.Interfaceを変える
	return types.Universe.Lookup("error").Type().Underlying().(*types.Interface), nil
}
