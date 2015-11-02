//TODO:
//1) Actually look after file, confirm that the file exist instead of giving panic error message.
//2) Investigate how Go tries to recover error-prone source code.

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/parser"
	"os"
	"path/filepath"
	"io/ioutil"
)

func main() {
	fmt.Println("This is the Go commom mistake detector... Looking for mistakes...\n")

	srcFile, err := getFilenameFromCommandLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	src, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error:\n")
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	//ast.Print(fset, f) //Print AST
	visitor := visitor{fileSet:fset}
	ast.Walk(&visitor, f)
}

type visitor struct {
	fileSet   *token.FileSet
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.GoStmt:
		if findRaceInGoRoutine(t) {
			fmt.Printf("$$ Warning: Potential race-condition found on line %d $$\n", getLineNumberInSourceCode(v.fileSet, t.Pos()))
		}
	case *ast.BadDecl:
		fmt.Println("#### Eroro in source code #####")
		os.Exit(1)
	case *ast.BadExpr:
		fmt.Println("#### Eroro in source code #####")
		os.Exit(1)
	case *ast.BadStmt:
		fmt.Println("#### Eroro in source code #####")
		os.Exit(1)
	}
	return v
}

func findRaceInGoRoutine(goNode *ast.GoStmt) (bool) {
	switch t := goNode.Call.Fun.(type) {
	case *ast.FuncLit:
		params := t.Type.Params.List
		for _, each := range t.Body.List {
			switch t1 := each.(type) {
			case *ast.ExprStmt:
				if (!validateParams(t1, params)) {
					return true;
				}
			}
		}
	}
	return false;
}

func validateParams(node *ast.ExprStmt, List []*ast.Field) (bool) {
	switch t := node.X.(type) {
	case *ast.CallExpr:
		for _, each := range t.Args {
			switch t1 := each.(type) {
			case *ast.Ident:
				if !containsListParam(t1, List) {
					return false;
				}
			}

		}
	}
	return true;
}

func containsListParam(ident *ast.Ident, List []*ast.Field) (found bool) {
	for _, each := range List {
		for _, each1 := range each.Names {
			if each1.Name == ident.Name {
				return true;
			}
		}
	}
	return false
}

func getLineNumberInSourceCode(fileSet *token.FileSet, position token.Pos) (int) {
	tokenFile := fileSet.File(position)
	return tokenFile.Line(position)
}

func getFilenameFromCommandLine() (srcFilename string, err error) {
	if len(os.Args) > 2 && os.Args[1] == "-s" {
		return os.Args[2], nil
	}
	err = fmt.Errorf("Usage: %s -s <go_source.go>\n", filepath.Base(os.Args[0]))
	return "", err
}
