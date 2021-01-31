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

func Chain(j Judge, mw ...Middleware) Middleware {
	fmt.Println("length:", len(mw))
	return func(next Mul) Mul {
		for i := len(mw) - 1; i >= 0; i-- {
			fmt.Println(i, j.judge[i], "isTrue?")
			if j.judge[i] {
				next = mw[i](next)
			} else {
				continue
			}
		}
		return next
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
	return func(next Mul) Mul {
		sum := ad.a + ad.b
		fmt.Println("sum:", sum)
		if sum > 10 {
			return nil
		}
		return func(ctx context.Context, c interface{}) (interface{}, error) {
			return next(ctx, c)
		}
	}
}

func log(l string) Middleware {
	return func(next Mul) Mul {
		fmt.Println("messages:", l)
		return func(ctx context.Context, c interface{}) (interface{}, error) {
			return next(ctx, c)
		}
	}
}

func main() {
	var j Judge
	j.judge = append(j.judge, true, false, true)
	a := &addNumber{a: 1, b: 3}
	first := Chain(j, log("this is a messages for mul."), log("111"), add(a))(mul)
	r, err := first(context.Background(), &mulNumber{9, 3})
	if err != nil {
		panic(err)
	}
	fmt.Println("mul result:", r)
	fmt.Println("--------------------------------------")

	second := Chain(j, log("this is a messages for sub."), log("222"), add(a))(sub)
	u, err := second(context.Background(), &subNumber{3, 2})
	if err != nil {
		panic(err)
	}
	fmt.Println("sub result:", u)
}