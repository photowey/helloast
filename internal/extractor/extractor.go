package extractor

import (
	"go/ast"
	"strings"

	"github.com/photowey/helloast/internal/astx"
)

var _ Extractor = (*extractor)(nil)

type Extractor interface {
	Package(_astx_ *astx.Ast)
	ValueSpec(_astx_ *astx.Ast)
	StructSpec(_astx_ *astx.Ast)
	InterfaceSpec(_astx_ *astx.Ast)
	FuncSpec(_astx_ *astx.Ast)
}

type extractor struct{}

func (e *extractor) Package(_astx_ *astx.Ast) {
	e.ValueSpec(_astx_)
	e.StructSpec(_astx_)
	e.InterfaceSpec(_astx_)
	e.FuncSpec(_astx_)
}

func (e *extractor) ValueSpec(_astx_ *astx.Ast) {
	for _, d := range _astx_.Ast.Decls {
		if strings.HasSuffix(_astx_.Path, "ping_handler.go") {
			switch decl := d.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch specVal := spec.(type) {
					case *ast.ValueSpec:
						values := make([]astx.Value, 0, len(specVal.Values))
						names := specVal.Names
						composite := false
						for idx, value := range specVal.Values {
							if bl, ok := value.(*ast.BasicLit); ok {
								values = append(values, astx.Value{
									Composite: false,
									Kind:      bl.Kind,
									Name:      names[idx].Name,
									Value:     bl.Value,
								})
							}
							if cl, ok := value.(*ast.CompositeLit); ok {
								for _, elt := range cl.Elts {
									if bl, okk := elt.(*ast.BasicLit); okk {
										composite = true
										values = append(values, astx.Value{
											Composite: composite,
											Kind:      bl.Kind,
											Name:      names[0].Name,
											Value:     bl.Value,
										})
									}
								}
							}
						}
						vs := &astx.ValueSpec{
							Pkg:       _astx_.Pkg,
							Tok:       decl.Tok,
							Composite: composite,
							Values:    values,
						}
						_astx_.Values = append(_astx_.Values, vs)
					}
				}
			}
		}
	}
}

func (e *extractor) StructSpec(_astx_ *astx.Ast) {
	for _, d := range _astx_.Ast.Decls {
		if strings.HasSuffix(_astx_.Path, "ping_handler.go") {
			switch decl := d.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch specVal := spec.(type) {
					case *ast.TypeSpec:
						if st, ok := specVal.Type.(*ast.StructType); ok {
							comments := make([]string, 0, len(decl.Doc.List))
							if decl.Doc != nil {
								for _, comment := range decl.Doc.List {
									comments = append(comments, comment.Text)
								}
							}

							ss := &astx.StructSpec{
								Pkg:         _astx_.Pkg,
								Name:        specVal.Name.String(),
								Comments:    comments,
								Fields:      make([]*astx.FieldSpec, 0),
								Methods:     make([]*astx.MethodSpec, 0),
								Annotations: make([]*astx.Annotation, 0),
							}

							ss.Type = st.Struct
							_astx_.Structs = append(_astx_.Structs, ss)
						}
					}
				}
			}
		}
	}
}

func (e *extractor) InterfaceSpec(_astx_ *astx.Ast) {
	for _, d := range _astx_.Ast.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch specVal := spec.(type) {
				case *ast.TypeSpec:
					if it, ok := specVal.Type.(*ast.InterfaceType); ok {
						comments := make([]string, 0, len(decl.Doc.List))
						if decl.Doc != nil {
							for _, comment := range decl.Doc.List {
								comments = append(comments, comment.Text)
							}
						}

						is := &astx.InterfaceSpec{
							Pkg:         _astx_.Pkg,
							Name:        specVal.Name.String(),
							Comments:    comments,
							Methods:     make([]*astx.MethodSpec, 0),
							Annotations: make([]*astx.Annotation, 0),
						}
						is.Type = it.Interface
						_astx_.Interfaces = append(_astx_.Interfaces, is)
					}
				}
			}
		}
	}
}

func (e *extractor) FuncSpec(_astx_ *astx.Ast) {
	for _, d := range _astx_.Ast.Decls {
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
					for _, spec := range _astx_.Structs {
						structName := spec.Name
						if structName == id.Name {
							ms := &astx.MethodSpec{
								Pkg:      _astx_.Pkg,
								Struct:   structName,
								Name:     decl.Name.String(),
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
					Pkg:      _astx_.Pkg,
					Name:     decl.Name.String(),
					Comments: comments,
					Params:   make([]*astx.ParamSpec, 0),
					Returns:  make([]*astx.ReturnSpec, 0),
				}
				_astx_.Funcs = append(_astx_.Funcs, fs)
			}
		}
	}
}

func NewExtractor() Extractor {
	return &extractor{}
}
