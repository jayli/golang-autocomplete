package main

import (
	"log"
	"reflect"
)

// kkk

func funcName() (string,int) {
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
