package test_pkg

import (
	"log"
	"reflect"
)

// kkk

const global_var string = "xxxxx"
const local_var string = "yyyy"

type Global_type_struct struct {
	Logf               string
	Context            []string
	ContextAAA         []string
	ContextBBB         []string
	Builtin            bool
	IgnoreCase         bool
	UnimportedPackages bool
}

var (
	g_is_server = 1
	g_cache     = 2
	g_format    = 3
	g_input     = 4
)

func funcName() (string, int) {
	const inner_global_var string = "gzzzzz"

	var aaaaaaa string
	return aaaaaaa, 1
}

func main() {

	var client_path string
	client_path = "asdf"
	log.Println(client_path)

	var a interface{}
	a = map[string]string{"a": "b"}
	log.Println(a.(map[string]string))
	log.Println(reflect.TypeOf(a))

	// if aaaa :=

}
