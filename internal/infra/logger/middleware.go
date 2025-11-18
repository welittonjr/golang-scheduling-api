package logger

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.NewString()
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Trace-ID", traceID)

		next.ServeHTTP(w, r)
	})
}
