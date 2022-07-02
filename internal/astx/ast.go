package astx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"github.com/photowey/helloast/internal/astx/loader"
)

const commentPrefix = string("//")

func NewAst(path string) *Ast {
	name := filepath.Base(path)
	_ast_, err := BuildAstFile(path)
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

func NewAstx(path string, lpkg *loader.Package) *Astx {
	pkg := lpkg.PkgPath
	name := filepath.Base(path)
	af, err := BuildAstFile(path)
	if err != nil {
		panic(err)
	}
	astx := &Astx{
		Package: lpkg,
		Path:    path,
		Name:    name,
		Pkg:     pkg,
		Ast:     af,
	}

	return astx
}

func BuildAstFile(path string) (*ast.File, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(path)
	fset := token.NewFileSet()

	return parser.ParseFile(fset, fileName, string(bytes), parser.ParseComments)
}
