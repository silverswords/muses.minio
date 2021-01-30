package main

import "fmt"

type number struct {
	a, b int
}

func add1(n *number) int {
	return n.a + n.b
}

type multiplier struct {
	n number
	c int
}

func mul1(m *multiplier) int {
	sum := add1(&m.n)
	product := m.c * sum
	return product
}

type logger struct {
	m multiplier
}

func log1(l *logger) {
	fmt.Println("result:", mul1(&l.m))
}

func main() {
	var l = logger{multiplier{number{2, 3}, 5}}
	log1(&l)
}
