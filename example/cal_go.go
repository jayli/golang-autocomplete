// kkkkkk
package main

import "fmt"
import "time"

func funcName() error {

	var count2 string
	return count2

}

type NewFunction struct {
	a string
}

func (nf NewFunction) doSth() {

}

func main(nf NewFunction) {
	tt := nf
	t1 := time.Now()
	count := int64(0)
	max := int64(9000000000)
	for i := int64(0); i < max; i++ {
		count += i
	}
	t2 := time.Now()
	fmt.Printf("cost:%d,count:%d\n", t2.Sub(t1)/1000000000, count)

}
