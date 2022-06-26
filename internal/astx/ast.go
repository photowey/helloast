package astx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
)

const commentPrefix = string("//")

type Ast struct {
	Path       string
	Name       string
	Pkg        string
	Values     []*ValueSpec
	Structs    []*StructSpec
	Interfaces []*InterfaceSpec
	Funcs      []*FuncSpec
	Ast        *ast.File
}

func NewAst(path string) *Ast {
	name := filepath.Base(path)
	_ast_, err := buildAstFile(path)
	if err != nil {
		panic(err)
	}
	astx := &Ast{
		Path:       path,
		Name:       name,
		Values:     make([]*ValueSpec, 0),
		Structs:    make([]*StructSpec, 0),
		Interfaces: make([]*InterfaceSpec, 0),
		Funcs:      make([]*FuncSpec, 0),
		Ast:        _ast_,
	}

	astx.Pkg = astx.Ast.Name.String()

	return astx
}

type AnalysisResult struct {
	PkgName     string
	RecvMethods map[string][]MethodInfo
	Funcs       []FuncInfo
}

type MethodInfo struct {
	PkgName    string
	RecvName   string
	MethodName string
	Comment    []string
}

type FuncInfo struct {
	PkgName  string
	FuncName string
	Comment  []string
}

func buildAstFile(path string) (*ast.File, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(path)
	fset := token.NewFileSet()

	return parser.ParseFile(fset, fileName, string(bytes), parser.ParseComments)
}
