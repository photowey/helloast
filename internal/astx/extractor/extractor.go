package extractor

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/photowey/helloast/internal/astx"
	"github.com/photowey/helloast/internal/astx/loader"
	"github.com/photowey/helloast/internal/astx/parser"
)

var _ Extractor = (*extractor)(nil)

type Extractor interface {
	Package(_astx_ *astx.Ast)
	ValueSpec(_astx_ *astx.Ast)
	StructSpec(_astx_ *astx.Ast)
	InterfaceSpec(_astx_ *astx.Ast)
	FuncSpec(_astx_ *astx.Ast)
	FieldSpec(_astx_ *astx.Ast)
	ScanPackage(rootPaths ...string) []*astx.AstSpec
}

type extractor struct{}

func (e *extractor) Package(_astx_ *astx.Ast) {
	e.ValueSpec(_astx_)
	e.StructSpec(_astx_)
	e.InterfaceSpec(_astx_)
	// e.FuncSpec(_astx_)
	fmt.Println("")
}

func (e *extractor) ValueSpec(_astx_ *astx.Ast) {
	for _, d := range _astx_.Ast.Decls {
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

							if fields := st.Fields; fields != nil && fields.List != nil {
								for _, field := range fields.List {
									fs := &astx.FieldSpec{
										Struct: specVal.Name.String(),
										Name:   field.Names[0].Name,
										Tags:   make([]*astx.TagSpec, 0),
										Ptr:    false,
									}

									// handle field's type
									switch expr := field.Type.(type) {
									case *ast.Ident: // Xxx Yyy `k:"v"`
										fs.Type = expr.Name
									case *ast.StarExpr: // Xxx *Yyy `k:"v"`
										switch starExpr := expr.X.(type) {
										case *ast.SelectorExpr:
											if sx, oke := starExpr.X.(*ast.Ident); oke {
												fs.Type = fmt.Sprintf("*%s.%s", sx.Name, starExpr.Sel.Name)
											} else {
												fs.Type = fmt.Sprintf("*%s", starExpr.Sel.Name)
											}
										case *ast.Ident:
											fs.Type = fmt.Sprintf("*%s", starExpr.Name)
										}

										fs.Ptr = true
									case *ast.SelectorExpr: //  Xxx yyy.Zzz `k:"v"`
										if sx, oke := expr.X.(*ast.Ident); oke {
											fs.Type = fmt.Sprintf("%s.%s", sx.Name, expr.Sel.Name)
										} else {
											fs.Type = expr.Sel.Name
										}
									}

									// handle field's tag
									if fieldTag := field.Tag; fieldTag != nil {
										tagValue := fieldTag.Value               // `xxx:"xv" yyy:"yv"`
										tagValue = tagValue[1 : len(tagValue)-1] // xxx:"xv" yyy:"yv"
										ts := &astx.TagSpec{
											Field: fs.Name,
											Tags:  make([]*astx.Tag, 0),
										}
										tvs := strings.Split(tagValue, " ")
										for _, tv := range tvs {
											kvs := strings.Split(tv, ":")
											k := kvs[0]                    // xxx | yyy
											v := kvs[1][1 : len(kvs[1])-1] // xv | yv
											tag := &astx.Tag{
												Name:  fieldTag.Value,
												Key:   k,
												Value: v,
											}
											ts.Tags = append(ts.Tags, tag)
										}

										fs.Tags = append(fs.Tags, ts)
									}

									ss.Fields = append(ss.Fields, fs)
								}
							}

							ss.Type = st.Struct
							_astx_.Structs = append(_astx_.Structs, ss)
						}
					}
				}
			}
		}
	}
	e.FuncSpec(_astx_)
}

func (e *extractor) InterfaceSpec(_astx_ *astx.Ast) {
	for _, d := range _astx_.Ast.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch specVal := spec.(type) {
				case *ast.TypeSpec:
					if it, ok := specVal.Type.(*ast.InterfaceType); ok {
						comments := make([]string, 0)
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
			comments := make([]string, 0)
			if decl.Doc != nil {
				for _, comment := range decl.Doc.List {
					comments = append(comments, comment.Text)
				}
			}
			if decl.Recv != nil {
				for _, field := range decl.Recv.List {
					switch ft := field.Type.(type) {
					case *ast.Ident:
						for _, spec := range _astx_.Structs {
							structName := spec.Name
							if structName == ft.Name {
								ms := &astx.MethodSpec{
									Pkg:      _astx_.Pkg,
									Struct:   structName,
									Name:     decl.Name.String(),
									Comments: comments,
									Params:   make([]*astx.ParamSpec, 0),
									Returns:  make([]*astx.ReturnSpec, 0),
								}

								hasParams := decl.Type != nil && decl.Type.Params != nil && decl.Type.Params.List != nil
								if hasParams {
									for _, param := range decl.Type.Params.List {
										for _, pn := range param.Names {
											pms := &astx.ParamSpec{
												Pkg:      _astx_.Pkg,
												FuncName: decl.Name.String(),
												Name:     pn.Name,
												Ptr:      false,
											}

											switch expr := param.Type.(type) {
											case *ast.Ident:
												pms.Type = expr.Name
											case *ast.StarExpr:
												switch starExpr := expr.X.(type) {
												case *ast.Ident:
													pms.Type = starExpr.Name
												case *ast.SelectorExpr:
													if x, okx := starExpr.X.(*ast.Ident); okx {
														pt := fmt.Sprintf("*%s.%s", x.Name, starExpr.Sel.Name)
														pms.Type = pt
													} else {
														pt := fmt.Sprintf("*%s", starExpr.Sel.Name)
														pms.Type = pt
													}
												}
												pms.Ptr = true
											case *ast.SelectorExpr:
												if x, okx := expr.X.(*ast.Ident); okx {
													pms.Type = fmt.Sprintf("%s.%s", x.Name, expr.Sel.Name)
												} else {
													pms.Type = expr.Sel.Name
												}
											}

											ms.Params = append(ms.Params, pms)
										}
									}
								}

								hasResults := decl.Type != nil && decl.Type.Results != nil && decl.Type.Results.List != nil
								if hasResults {
									for _, rvt := range decl.Type.Results.List {
										rs := &astx.ReturnSpec{
											Pkg:      _astx_.Pkg,
											FuncName: decl.Name.String(),
											Ptr:      false,
										}
										if names := rvt.Names; names != nil {
											rs.Name = rvt.Names[0].Name
										}

										switch rvtType := rvt.Type.(type) {
										case *ast.Ident:
											rs.Type = rvtType.Name
										case *ast.StarExpr:
											switch xt := rvtType.X.(type) {
											case *ast.Ident:
												rs.Type = fmt.Sprintf("*%s", xt.Name)
											case *ast.SelectorExpr:
												if x, okx := xt.X.(*ast.Ident); okx {
													rs.Type = fmt.Sprintf("*%s.%s", x.Name, xt.Sel.Name)
												} else {
													rs.Type = fmt.Sprintf("*%s", xt.Sel.Name)
												}
											}

											rs.Ptr = true
										}

										ms.Returns = append(ms.Returns, rs)
									}
								}
								spec.Methods = append(spec.Methods, ms)
							}
						}
					case *ast.StarExpr:
						// do nothing
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

func (e *extractor) FieldSpec(_astx_ *astx.Ast) {
}

func (e *extractor) ScanPackage(rootPaths ...string) []*astx.AstSpec {
	roots, err := loader.LoadRoots(rootPaths...)
	if err != nil {
		panic(err)
	}

	ass := make([]*astx.AstSpec, 0, len(roots))

	for _, root := range roots {
		as := parser.Parse(root)
		ass = append(ass, as)
	}

	return ass
}

func NewExtractor() Extractor {
	return &extractor{}
}
