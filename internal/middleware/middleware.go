package middleware

import "net/http"

type Middleware func(Handler) Handler

type Handler func(http.ResponseWriter, *http.Request)

func ApplyMiddleware(main Handler, m ...Middleware) Handler {
	// base case
	if len(m) == 0 {
		return main
	}

	// recursion
	return m[0](ApplyMiddleware(main, m[1:cap(m)]...))
}
