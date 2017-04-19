package mw

import (
	"net/http"
)

type responseWriterKey int

type handlerAdapter struct {
	sub Handler
}

func (adapter *handlerAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := NewResponse(w)
	adapter.sub(resp, r)
	// Ignore error here I guess, although I would probably reconsider that.
}

// IntoHTTPMiddleware converts go-mw middleware into a net/http middleware, with
// one caveat: The Response input should be treated just like a ResponseWriter,
// i.e. modification of headers after it's handed downstream will (most likely)
// have no effect.
func IntoHTTPMiddleware(m Middleware) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		handler := func(resp *Response, r *http.Request) error {
			h.ServeHTTP(resp.Writer, r)
			return nil
		}
		return &handlerAdapter{sub: m(handler)}
	}
}
