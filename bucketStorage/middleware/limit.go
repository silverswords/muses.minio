package middleware

import (
	"context"
)

type Condition struct {
	ObjectSize uint64
	AllSize uint64
	Threshold uint64
}

func Limit(c *Condition) Middleware {
	return func(next Handler) Handler {
		c.AllSize += c.ObjectSize
		if c.AllSize > c.Threshold {
			return nil
		}
		return func(ctx context.Context, p interface{}) (interface{}, error) {
			return next(ctx, p)
		}
	}
}
