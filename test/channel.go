package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5*time.Second)
		close(c)
	}()

	x, ok := <-c
	fmt.Println("channel", x, ok)
	select {
	case <-c:
		fmt.Printf("uuu %v.\n", time.Since(start))
	}
}
