package middleware

import (
	"github.com/silverswords/muses.minio/mos"
)

type middleware []interface{}

type Bucket mos.Bucket

func (b *Bucket) RegisterMiddlewares(middlewares ...middleware) {
	for _, mid := range middlewares {
		b.Middlewares = append(b.Middlewares, mid)
	}
}