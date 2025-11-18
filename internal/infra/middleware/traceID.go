package middleware

import (
	http "scheduling/internal/infra/gin"

	"github.com/google/uuid"
)

const (
	traceIDKey string = "trace_id"
)

func TraceIDMiddleware() http.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) error {
			traceID := uuid.New().String()

			ctx.Set(traceIDKey, traceID)
			ctx.Header("X-Trace-ID", traceID)
			
			return next(ctx)
		}
	}
}