package middleware

import (
	"net/http"

	"github.com/disiqueira/PoContext/context"
)

func TraceID() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("ot-tracer-traceid")
			ctx := context.WithTraceID(r.Context(), traceID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
