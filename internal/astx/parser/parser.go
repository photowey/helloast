package parser

import (
	"github.com/photowey/helloast/internal/astx"
	"github.com/photowey/helloast/internal/astx/loader"
)

func Parse(pkg *loader.Package) *astx.AstSpec {
	return &astx.AstSpec{
		Path: pkg.PkgPath,
		Pkgs: make([]*astx.PackageSpec, 0),
	}
}
