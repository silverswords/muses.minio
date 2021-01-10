package mos

import (
	"github.com/silverswords/muses.minio/mos/driver"
	"sync"
)

type Bucket struct {
	b      driver.Bucket

	// mu protects the closed variable.
	// Read locks are kept to allow holding a read lock for long-running calls,
	// and thereby prevent closing until a call finishes.
	mu          sync.RWMutex
	closed      bool
	Middlewares []interface{}
}

var NewBucket = newBucket

func newBucket(b driver.Bucket) *Bucket {
	return &Bucket{
		b: b,
	}
}
