// go run client.go -cursor=11
// https://godoc.org/go/ast

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var g_input = flag.String("in", "", "输入文件")
var g_cursor = flag.String("cursor", "", "游标位置")

var file []byte
var filename string
var cursor int64

// 本地调试
const debugger bool = true

func funcName1() {

	var aaaaaaa int
	var bbbbbbb int
	aaaaaaa = 1
	bbbbbbb = 2

	log.Println(aaaaaaa)
	log.Println(bbbbbbb)
	return

}

// 这里要处理以路径方式引入进来的 pkgs
// 比如 [go/ast go/parser go/token go/types]
func getGlobalImports(f *ast.File) []string {
	var imports_map []string
	imports_map = make([]string, 0)

	for _, v := range f.Imports {
		imports_map = append(imports_map, strings.Replace(v.Path.Value, "\"", "", -1))
	}
	return imports_map
}

func getGlobalPkgs(f *ast.File) []string {
	var pkgs []string
	pkgs = make([]string, 0)
	return append(pkgs, f.Name.Name)
}

type GlobalIdent struct {
	consts []string
	vars   []string
	types  []string
}

// 返回一个 Spec 中的  常量列表，变量列表，类型列表
func getValueSpecNames(vsp *ast.Spec) *GlobalIdent {

	global_ident := &GlobalIdent{
		consts: make([]string, 0),
		vars:   make([]string, 0),
		types:  make([]string, 0),
	}

	switch (*vsp).(type) {
	case *ast.ValueSpec:
		local_vsp := (*vsp).(*ast.ValueSpec)
		for _, id := range local_vsp.Names {
			var kind = id.Obj.Kind.String()
			var name = id.Obj.Name

			switch kind {
			case "const":
				global_ident.consts = append(global_ident.consts, name)
			case "var":
				global_ident.vars = append(global_ident.vars, name)
			default:
				// Do Nothing
			}
		}
	case *ast.TypeSpec:
		// type 语法不会像 const 和 var 一样用括号一次定义一堆
		local_vsp := (*vsp).(*ast.TypeSpec)
		global_ident.types = append(global_ident.types, local_vsp.Name.Name)
	default:
		// Do Nothing
	}

	return global_ident
}

func getGlobalTypes(f *ast.File) []string {
	var types []string
	types = make([]string, 0)

	for _, decl := range f.Decls {
		if _, ok := decl.(*ast.GenDecl); ok == false {
			continue
		}

		if fmt.Sprint(decl.(*ast.GenDecl).Tok) == "type" {
			for _, vsp := range decl.(*ast.GenDecl).Specs {
				tmp_types := getValueSpecNames(&vsp).types
				types = append(types, tmp_types...)
			}
		}

	}
	return types
}

// getVars
func getGlobalVars(f *ast.File) []string {

	var vars []string
	vars = make([]string, 0)

	for _, decl := range f.Decls {
		if _, ok := decl.(*ast.GenDecl); ok == false {
			continue
		}
		var field reflect.Value = reflect.ValueOf(decl).Elem().FieldByName("Tok")
		if field.IsValid() && fmt.Sprint(decl.(*ast.GenDecl).Tok) == "var" {
			for _, vsp := range decl.(*ast.GenDecl).Specs {
				tmp_vars := getValueSpecNames(&vsp).vars
				vars = append(vars, tmp_vars...)
			}
		}
	}
	return vars
}

// getConsts
func getGlobalConsts(f *ast.File) []string {

	var consts []string
	consts = make([]string, 0)

	for _, decl := range f.Decls {
		// 类型断言
		if _, ok := decl.(*ast.GenDecl); ok == false {
			continue
		}

		// 判断 struct 成员是否合法
		var field reflect.Value = reflect.ValueOf(decl).Elem().FieldByName("Tok")

		if field.IsValid() && fmt.Sprint(decl.(*ast.GenDecl).Tok) == "const" {
			for _, vsp := range decl.(*ast.GenDecl).Specs {
				tmp_consts := getValueSpecNames(&vsp).consts
				consts = append(consts, tmp_consts...)
			}
		}
	}
	return consts
}

func getGlobalFuncs(f *ast.File) []string {
	var funcs []string
	funcs = make([]string, 0)

	for _, decl := range f.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok == false {
			continue
		}
		funcs = append(funcs, decl.(*ast.FuncDecl).Name.Name)
	}

	return funcs
}

func LogRootMembers(f []ast.Decl) {
	for _, v := range f {
		log.Println(reflect.ValueOf(v))
	}
}

func main() {
	flag.Parse()

	if debugger == true {
		*g_input = "example/test.go"
	}

	file, _ = ioutil.ReadFile(*g_input)
	filename, _ = filepath.Abs(*g_input)
	cursor, _ = strconv.ParseInt(*g_cursor, 10, 0)

	if debugger == true {
		cursor = 497
	}

	log.Println("------------[[[---------------")
	log.Println("------------[[[---------------")
	cc, a, b := deduceCursorContext(file, int(cursor))
	fmt.Print(string(file))
	log.Println(">>>   ", cc)
	log.Println(">>>   ", a)
	log.Println(">>>   ", b)
	log.Println("------------]]]---------------")
	log.Println("------------]]]---------------")

	//source_file := "/Users/bachi/jayli/golang-autocomplete/example/cal_go.go"
	var source_file string
	source_file = "/Users/bachi/jayli/golang-autocomplete/gocode/utils.go"
	source_file = "/Users/bachi/jayli/golang-autocomplete/gocode/internal/suggest/candidate.go"
	source_file = "/Users/bachi/jayli/golang-autocomplete/gocode/internal/suggest/suggest.go"
	source_file = "/Users/bachi/jayli/golang-autocomplete/example/test.go"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, source_file, nil, parser.AllErrors)

	if err != nil {
		panic(err)
	}

	log.Println("--")
	// ast.Print(fset, f.Decls)
	log.Println("--")
	LogRootMembers(f.Decls)

	log.Println(">>getImports: \t", getGlobalImports(f))
	log.Println(">>getPkgs: \t\t", getGlobalPkgs(f))
	log.Println(">>getConsts: \t", getGlobalConsts(f))
	log.Println(">>getVars: \t\t", getGlobalVars(f))
	log.Println(">>getTypes: \t", getGlobalTypes(f))
	log.Println(">>getFuncs: \t", getGlobalFuncs(f))

	log.Println("------------EOF-A---------------")

	// ast.Inspect(f, func(n ast.Node) bool {
	// 	var s string
	// 	switch x := n.(type) {
	// 	case *ast.BasicLit:
	// 		s = x.Value
	// 	case *ast.Ident:
	// 		s = x.Name
	// 	}

	// 	if s != "" {
	// 		// fmt.Printf("\t%s\n", s)
	// 	}
	// 	return true
	// })

	// res := gocodeAutoComplete(filename, file, cursor)

	log.Println("------------EOF-B---------------")

	os.Exit(1)
}
