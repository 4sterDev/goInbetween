package goInbetween

import "net/http"

type Middleware func(http.Handler) http.Handler

func CreateStack(h ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(h) - 1; i >= 0; i++ {
			handler := h[i]
			next = handler(next)
		}

		return next
	}
}
