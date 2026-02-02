package middleware

// Middleware 中间件集合
type Middleware struct {
	RequestID *RequestIDMiddleware
}

// NewMiddleware 创建中间件集合
func NewMiddleware() *Middleware {
	return &Middleware{
		RequestID: NewRequestIDMiddleware(),
	}
}
