package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/photowey/helloast/internal/astx"
	"github.com/photowey/helloast/internal/astx/extractor"
	"github.com/photowey/helloast/internal/astx/loader"
)

const (
	FileSuffix     = ".go"
	App            = "@App"
	RestController = "@RestController"
)

var _extractor_ = extractor.NewExtractor()

func Run() error {
	return run()
}

func run() error {
	wd, _ := os.Getwd()
	fmt.Printf("App.start on:%s \n", wd)

	// rf := filepath.Join(wd, "./...")

	doScan(wd)
	// loadRoots(rf)
	// scanPkg(rf)

	return nil
}

func doScan(scanPath string) {
	scanPath = filepath.Clean(scanPath)

	_ = filepath.Walk(scanPath, walk)
}

func walk(path string, info os.FileInfo, err error) error {
	if info == nil {
		return err
	}

	if isFile(info) && isGoFile(info) {
		scan(path)
	}

	return nil
}

func isFile(info os.FileInfo) bool {
	return !info.IsDir()
}

func isGoFile(info os.FileInfo) bool {
	return strings.HasSuffix(info.Name(), FileSuffix)
}

func scan(path string) {
	ast := astx.NewAst(path)
	_extractor_.Package(ast)
}

func loadRoots(rootPaths ...string) {
	roots, err := loader.LoadRoots(rootPaths...)
	if err != nil {
		panic(err)
	}

	for _, root := range roots {
		fmt.Println(root)
	}
}

func scanPkg(rootPaths ...string) {
	_extractor_.ScanPackage(rootPaths...)
}
