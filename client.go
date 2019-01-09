// go run client.go -cursor=11

package main

import (
	"./internal/suggest"
	"flag"
	"io/ioutil"
	"log"
	//"os"
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

	/*
		Candidates := []suggest.Candidate{
			{Class: "PANIC", Name: "PANIC", Type: "PANIC"},
		}
	*/

	cfg := suggest.Config{
		Context:            Context,
		Builtin:            false,
		IgnoreCase:         true,
		UnimportedPackages: false,
		Logf:               func(string, ...interface{}) {},
	}

	candidates, d := cfg.Suggest(filename, file, int(cursor))
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
		cursor = 269
	}

	log.Println("file: ", file)
	log.Println("filename: ", filename)
	log.Println("cursor: ", cursor)

	log.Println("-----------------------------")

	res := gocodeAutoComplete(filename, file, cursor)

	log.Println(res.Candidates)
	log.Println(res.Len)

	log.Println("-----------------------------")

	// prepareFilenameDataCursor()
}
