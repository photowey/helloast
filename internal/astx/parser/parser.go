package parser

import (
	"go/ast"

	"github.com/photowey/helloast/internal/astx"
	"github.com/photowey/helloast/internal/astx/loader"
)

var (
	_        Parser = (*parser)(nil)
	_parser_ Parser
)

func init() {
	_parser_ = NewParser()
}

type Parser interface {
	Parse(pkg *loader.Package) *astx.AstSpec
	ParseStructs(aw *astx.Astx) *astx.PackageSpec
	ParseMethods(aw *astx.Astx, ps *astx.PackageSpec)
}

type parser struct{}

func (psr parser) Parse(pkg *loader.Package) *astx.AstSpec {
	pkgs := make([]*astx.PackageSpec, 0, len(pkg.CompiledGoFiles))
	for _, cf := range pkg.CompiledGoFiles {
		aw := astx.NewAstx(cf, pkg)
		if aw.Ast.Comments == nil {
			continue
		}

		ps := psr.ParseStructs(aw)
		if len(ps.Structs) == 0 {
			continue
		}

		psr.ParseMethods(aw, ps)

		pkgs = append(pkgs, ps)
	}

	return &astx.AstSpec{
		PkgPath: pkg.PkgPath,
		Pkgs:    make([]*astx.PackageSpec, 0),
	}
}

func (psr parser) ParseStructs(aw *astx.Astx) *astx.PackageSpec {
	ps := &astx.PackageSpec{
		Pkg:        aw.Package.PkgPath,
		Alias:      aw.Package.Name,
		Structs:    make([]*astx.StructSpec, 0),
		Interfaces: make([]*astx.InterfaceSpec, 0),
		Funcs:      make([]*astx.FuncSpec, 0),
	}
	for _, d := range aw.Ast.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
		SPEC:
			for _, spec := range decl.Specs {
				switch specVal := spec.(type) {
				case *ast.TypeSpec:
					if st, ok := specVal.Type.(*ast.StructType); ok {
						if decl.Doc == nil {
							continue SPEC
						}
						comments := make([]string, 0, len(decl.Doc.List))
						if decl.Doc != nil {
							for _, comment := range decl.Doc.List {
								comments = append(comments, comment.Text)
							}
						}

						ss := &astx.StructSpec{
							Pkg:         ps.Pkg,
							Name:        specVal.Name.String(),
							Comments:    comments,
							Fields:      make([]*astx.FieldSpec, 0),
							Methods:     make([]*astx.MethodSpec, 0),
							Annotations: make([]*astx.Annotation, 0),
						}

						ss.Type = st.Struct
						ps.Structs = append(ps.Structs, ss)
					}
				}
			}
		}
	}

	return ps
}

func (psr parser) ParseMethods(aw *astx.Astx, ps *astx.PackageSpec) {
	for _, d := range aw.Ast.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			comments := make([]string, 0, len(decl.Doc.List))
			if decl.Doc != nil {
				for _, comment := range decl.Doc.List {
					comments = append(comments, comment.Text)
				}
			}
			if decl.Recv != nil {
				for _, field := range decl.Recv.List {
					id := field.Type.(*ast.Ident)
					for _, spec := range ps.Structs {
						structName := spec.Name
						if structName == id.Name {
							ms := &astx.MethodSpec{
								Pkg:      aw.Pkg,
								Name:     decl.Name.String(),
								Struct:   structName,
								Comments: comments,
								Params:   make([]*astx.ParamSpec, 0),
								Returns:  make([]*astx.ReturnSpec, 0),
							}
							spec.Methods = append(spec.Methods, ms)
						}
					}
				}
			} else {
				fs := &astx.FuncSpec{
					Pkg:      aw.Pkg,
					Name:     decl.Name.String(),
					Comments: comments,
					Params:   make([]*astx.ParamSpec, 0),
					Returns:  make([]*astx.ReturnSpec, 0),
				}
				ps.Funcs = append(ps.Funcs, fs)
			}
		}
	}
}

func NewParser() Parser {
	return &parser{}
}

func Parse(pkg *loader.Package) *astx.AstSpec {
	return _parser_.Parse(pkg)
}
