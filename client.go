// go run client.go -cursor=11
// https://godoc.org/go/ast

package main

import (
	"./internal/suggest"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	ast.Print(fset, f)
	log.Println("--")

	log.Println(">>")
	log.Println("<<")

	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}

		if s != "" {
			// fmt.Printf("\t%s\n", s)
		}
		return true
	})

	// res := gocodeAutoComplete(filename, file, cursor)

	log.Println("------------EOF---------------")

	os.Exit(1)
}
