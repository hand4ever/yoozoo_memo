package main

import (
	"fmt"
)

type double float64


func (a double) IsEqual(b double) bool {
	fmt.Println("func1")
	var r = a - b
	if r == 0.0 {
		return true
	} else if r < 0.0 {
		return r > -0.0001
	}
	return r < 0.0001
}

func IsEqual(a, b float64) bool {
	fmt.Println("func2")
	var r = a - b
	if r == 0.0 {
		return true
	} else if r < 0.0 {
		return r > -0.0001
	}
	return r < 0.0001
}

// func ReadFull(r Reader, buf []byte) (n int, err error) {
// 	for len(buf) > 0 && err == nil {
// 		var nr int
// 		nr, err = r.Read(buf)
// 		n += nr
// 		buf = buf[nr:]
// 	}
// 	return
// }

func main() {
	var a double = 1.999999
	var b double = 1.9999998
	fmt.Println(a.IsEqual(b))
	fmt.Println(a.IsEqual(3))
	fmt.Println(IsEqual((float64)(a), (float64)(b)))

	fc := func(msg string) {
		fmt.Println("you say: ", msg)
	}
	fmt.Printf("%T \n", fc)
	fc("hello, my love")

	func(msg string) {
		fmt.Println("say : ", msg)
	}("I love to code")
}
