package main

import "fmt"

func main() {
	var buf []byte
	fmt.Println(buf)
	if buf == nil {
		fmt.Println("true")
	}
}
