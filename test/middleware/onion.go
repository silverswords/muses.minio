package main

import "fmt"

type number struct {
	a, b int
}

func add(n *number) int {
	return n.a + n.b
}

type multiplier struct {
	n number
	c int
}

func mul(m *multiplier) int {
	sum := add(&m.n)
	product := m.c * sum
	return product
}

type logger struct {
	m multiplier
}

func log(l *logger) {
	fmt.Println("result:", mul(&l.m))
}

func main() {
	var l = logger{multiplier{number{2, 3}, 5}}
	log(&l)
}
