package test_pkg

import (
	"log"
	"reflect"
)

// kkk

const global_var string = "xxxxx"
const local_var string = "yyyy"

func funcName() (string, int) {
	const inner_global_var string = "zzzzzz"

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

}
