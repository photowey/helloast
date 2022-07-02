package astx

import (
	"go/ast"
	"go/token"

	"github.com/photowey/helloast/internal/astx/loader"
)

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

type Astx struct {
	*loader.Package
	Path string
	Name string
	Pkg  string
	Ast  *ast.File
}

type AstSpec struct {
	ID      string
	Name    string
	PkgPath string
	Pkgs    []*PackageSpec
}

type ValueSpec struct {
	Pkg       string
	Tok       token.Token
	Composite bool
	Values    []Value
}

type Value struct {
	Composite bool
	Kind      token.Token
	Name      string
	Value     any
}

type PackageSpec struct {
	Pkg        string
	Alias      string
	Structs    []*StructSpec
	Interfaces []*InterfaceSpec
	Funcs      []*FuncSpec
}

type StructSpec struct {
	Pkg         string
	Alias       string
	Name        string
	Type        token.Pos
	Comments    []string
	Fields      []*FieldSpec
	Methods     []*MethodSpec
	Annotations []*Annotation
}

type FieldSpec struct {
	Struct string
	Name   string
	Type   string
	Tags   []TagSpec
}

type TagSpec struct {
	Field string
	Tags  []Tag
}
type Tag struct {
	Key   string
	Value string
}

type InterfaceSpec struct {
	Pkg         string
	Name        string
	Type        token.Pos
	Comments    []string
	Methods     []*MethodSpec
	Annotations []*Annotation
}

type Annotation struct {
	Pkg    string
	Name   string
	Values string // maybe json ?
}

type MethodSpec struct {
	Pkg      string
	Struct   string
	Name     string
	Comments []string
	Params   []*ParamSpec
	Returns  []*ReturnSpec
}

type FuncSpec struct {
	Pkg      string
	Name     string
	Comments []string
	Params   []*ParamSpec
	Returns  []*ReturnSpec
}

type ParamSpec struct {
	Pkg      string
	FuncName string
	Name     string
	Type     string
}

type ReturnSpec struct {
	Pkg      string
	FuncName string
	Name     string
	Type     string
}
