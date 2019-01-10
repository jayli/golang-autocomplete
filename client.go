// go run client.go -cursor=11

package main

import (
	"./internal/suggest"
	"flag"
	"io/ioutil"
	"log"
	//"os"
	"go/ast"
	"go/parser"
	"go/token"
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

	src := `

		// kkkkkk 
		package main

		import "fmt"
		import "time"

		func main() {
			t1 := time.Now()
			count := int64(0)
			max := int64(9000000000)
			for i := int64(0); i < max; i++ {
				count += i
			}
			t2 := time.Now()
			fmt.Printf("cost:%d,count:%d\n", t2.Sub(t1)/1000000000, count)
		}
	`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)

	if err != nil {
		panic(err)
	}

	cmap := ast.NewCommentMap(fset, f, f.Comments)

	_ = cmap
	log.Println(">>>", f)

	log.Println("-----------------------------")

	// res := gocodeAutoComplete(filename, file, cursor)

	// log.Println("res.Candidates:", res.Candidates)
	// log.Println("res.Len:", res.Len)

	// log.Println("-----------------------------")
	log.Println("------------EOF---------------")

	// prepareFilenameDataCursor()
}
