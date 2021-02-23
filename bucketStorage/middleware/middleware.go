package middleware

import "context"

type Judge map[string]bool

type Handler func(context.Context, interface{}) (interface{}, error)

type Middleware func(Handler) Handler

func Chain(j Judge, mw ...Middleware) Middleware {
	var arr = make([]bool, len(j))
	for _, value := range j {
		arr = append(arr, value)
	}
	return func(next Handler) Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			if arr[i] {
				next = mw[i](next)
			} else {
				continue
			}
		}
		return next
	}
}

