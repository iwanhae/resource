//go:build ci

package openapi3

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func ParseDoc(t reflect.Type) {
	pkgPath := t.PkgPath()
	pkgPath = strings.TrimSuffix(pkgPath, "_test")

	// Use go list to get the directory of the package
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkgPath)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running go list:", err)
		return
	}

	pkgDir := strings.TrimSpace(string(output))

	// Parse the package directory
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing package:", err)
		return
	}

	for _, pkg := range pkgs {
		docPkg := doc.New(pkg, pkgPath, 0)
		for _, typ := range docPkg.Types {
			fmt.Println("Type:", typ.Name)
			fmt.Println("Doc:", typ.Doc)
			for _, c := range typ.Consts {
				for _, name := range c.Names {
					fmt.Printf("Constant: %s\n", name)
				}
			}
		}
	}
}
