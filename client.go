// go run client.go -cursor=11
// https://godoc.org/go/ast

package main

import (
	"./internal/suggest"
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

type AutoCompleteReply struct {
	Candidates []suggest.Candidate
	Len        int
}

func gocodeAutoComplete(filename string, file []byte, cursor int64) *AutoCompleteReply {
	Context := &suggest.PackedContext{
		Dir: filepath.Dir(filename),
	}

	//log.Println("file: ", file)

	cfg := suggest.Config{
		Context:            Context,
		Builtin:            false,
		IgnoreCase:         false,
		UnimportedPackages: false,
		Logf:               func(string, ...interface{}) {},
	}

	// candidates, d := cfg.Suggest(filename, file, int(cursor))
	f := "/Users/bachi/jayli/golang-autocomplete/example/cal_go.go"
	data, _ := ioutil.ReadFile(f)
	candidates, d := cfg.Suggest(f, data, 268)

	log.Println("candidates: ", candidates)

	if candidates == nil {
		candidates = []suggest.Candidate{}
	}

	return &AutoCompleteReply{
		Candidates: candidates,
		Len:        d,
	}
}

func funcName1() {

	var aaaaaaa int
	var bbbbbbb int
	aaaaaaa = 1
	bbbbbbb = 2

	log.Println(aaaaaaa)
	log.Println(bbbbbbb)
	return

}

type ConstList struct {
	Candidates []suggest.Candidate
	Len        int
}

func getGlobalConst() []ConstList {

	return nil

}

func getImports(f *ast.File) []string {
	var imports_map []string
	imports_map = make([]string, 0)

	for _, v := range f.Imports {
		imports_map = append(imports_map, strings.Replace(v.Path.Value, "\"", "", -1))
	}
	return imports_map
}

func getPkgs(f *ast.File) []string {
	// log.Println(">>>>>> ", reflect.TypeOf(nil))
	var pkgs []string
	pkgs = make([]string, 0)
	log.Println(reflect.TypeOf(f.Name.Name))
	return append(pkgs, f.Name.Name)
}

// 返回一个 Spec 中的  常量列表，变量列表
func getValueSpecNames(vsp *ast.Spec) ([]string, []string) {
	var const_names []string
	var var_names []string
	const_names = make([]string, 0)
	var_names = make([]string, 0)

	for _, id := range (*vsp).(*ast.ValueSpec).Names {
		var kind = id.Obj.Kind.String()
		var name = id.Obj.Name

		switch kind {
		case "const":
			const_names = append(const_names, name)
		case "var":
			var_names = append(var_names, name)
		default:
			// Do Nothing
		}
	}
	return const_names, var_names
}

func getConsts(f *ast.File) []string {

	var consts []string
	consts = make([]string, 0)

	for _, decl := range f.Decls {
		// 类型断言
		if _, ok := decl.(*ast.GenDecl); ok == false {
			continue
		}

		/*
			// 另一种类型断言的方法
			if fmt.Sprint(reflect.TypeOf(decl)) != "*ast.GenDecl" {
				continue
			}
		*/

		// 判断 struct 成员是否合法
		var field reflect.Value = reflect.ValueOf(decl).Elem().FieldByName("Tok")

		if field.IsValid() && fmt.Sprint(decl.(*ast.GenDecl).Tok) == "const" {
			for _, vsp := range decl.(*ast.GenDecl).Specs {
				tmp_consts, _ := getValueSpecNames(&vsp)
				consts = append(consts, tmp_consts...)
			}
		}
	}
	return consts
}

func LogRootMembers(f []ast.Decl) {
	log.Println(reflect.TypeOf(f))
	for _, v := range f {
		log.Println(reflect.TypeOf(v))
	}
}

func main() {
	flag.Parse()

	if debugger == true {
		*g_input = "example/cal_go.go"
	}

	file, _ = ioutil.ReadFile(*g_input)
	filename, _ = filepath.Abs(*g_input)
	cursor, _ = strconv.ParseInt(*g_cursor, 10, 0)

	if debugger == true {
		cursor = 268
	}

	//source_file := "/Users/bachi/jayli/golang-autocomplete/example/cal_go.go"
	source_file := "/Users/bachi/jayli/golang-autocomplete/example/test.go"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, source_file, nil, parser.AllErrors)

	if err != nil {
		panic(err)
	}

	log.Println("--")
	ast.Print(fset, f.Decls)
	log.Println("--")
	LogRootMembers(f.Decls)

	log.Println(">>getImports: ", getImports(f))
	log.Println(">>getPkgs: ", getPkgs(f))
	log.Println(">>getConsts: ", getConsts(f))

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

	log.Println("------------EOF---------------")

	os.Exit(1)
}
