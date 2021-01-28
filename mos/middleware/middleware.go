package middleware

type Middleware interface {
	usable() bool
}

type LimitMiddleware struct {
	LimitUpload int64
	AllSize     int64
	Use         bool
}

func (mw *LimitMiddleware) usable() bool {
	return mw.Use
}

func Usable(mw Middleware) bool {
	return mw.usable()
}
