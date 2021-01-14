package main

func add(a, b int) int {
	c := a + b
	return c
}

func inc(a, b int) int {
	c := a - b
	return c
}

func main() {
	var a = 4
	var b = 3
	add(a, b)
	inc(a, b)
}
