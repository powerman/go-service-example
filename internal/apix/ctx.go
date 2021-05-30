package apix

import "context"

type Ctx = context.Context

type contextKey int

const (
	_ contextKey = iota
	contextKeyRemoteIP
)

// NewContextWithRemoteIP creates and returns new context containing given
// remoteIP.
func NewContextWithRemoteIP(ctx Ctx, remoteIP string) Ctx {
	return context.WithValue(ctx, contextKeyRemoteIP, remoteIP)
}

// FromContext returns values describing request stored in ctx, if any.
func FromContext(ctx Ctx) (remoteIP string) {
	remoteIP, _ = ctx.Value(contextKeyRemoteIP).(string)
	return
}
