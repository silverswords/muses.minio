package middleware

type Middlewares struct {
	AllSize       int64
	ResourceLimit func(bucketName string, objectSize int64) int64
}
