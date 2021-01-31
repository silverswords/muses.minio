package main

import (
	"context"
	"fmt"
)

type Judge struct {
	judge []bool
}

type Mul func(ctx context.Context, c interface{}) (response interface{}, err error)

type Middleware func(Mul) Mul

func Chain(j Judge, outer Middleware, others ...Middleware) Middleware {
	return func(next Mul) Mul {
		for i := len(others) - 1; i >= 0; i-- {
			if j.judge[i] {
				next = others[i](next)
			}
			break
		}
		return outer(next)
	}
}

type addNumber struct {
	a, b int
}

type mulNumber struct {
	sum, c int
}

type subNumber struct {
	a, b int
}

func sub(ctx context.Context, mn interface{}) (interface{}, error) {
	var diff int
	if m, ok := mn.(*subNumber); ok {
		diff = m.a - m.b
	}
	return diff, nil
}

func mul(ctx context.Context, mn interface{}) (interface{}, error) {
	var product int
	if m, ok := mn.(*mulNumber); ok {
		product = m.sum * m.c
	}
	return product, nil
}

func add(ad *addNumber) Middleware {
	sum := ad.a + ad.b
	if sum > 10 {
		return nil
	}
	return func(next Mul) Mul {
		return func(ctx context.Context, c interface{}) (interface{}, error) {
			return next(ctx, c)
		}
	}
}

func log(l string) Middleware {
	fmt.Println("messages:", l)
	return func(next Mul) Mul {
		return func(ctx context.Context, c interface{}) (interface{}, error) {
			return next(ctx, c)
		}
	}
}

func main() {
	var j Judge
	j.judge[0] = true
	j.judge[1] = true
	a := &addNumber{a: 1, b: 3}
	first := Chain(j, add(a), log("this is a messages for mul."))(mul)
	r, err := first(context.Background(), &mulNumber{9, 3})
	if err != nil {
		panic(err)
	}
	fmt.Println("mul result:", r)

	second := Chain(j, add(a), log("this is a messages for sub."))(sub)
	u, err := second(context.Background(), &subNumber{3, 2})
	if err != nil {
		panic(err)
	}
	fmt.Println("sub result:", u)
}