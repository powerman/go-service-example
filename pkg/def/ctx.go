package def

import "context"

type Ctx = context.Context

type ctxKey int

const keyRemoteIP ctxKey = 0

// NewContextWithRemoteIP creates and returns new context containing given
// remoteIP.
func NewContextWithRemoteIP(ctx Ctx, remoteIP string) Ctx {
	return context.WithValue(ctx, keyRemoteIP, remoteIP)
}

// FromContext returns all values which might be stored in context using
// NewContext* functions of this package.
func FromContext(ctx Ctx) (remoteIP string) {
	remoteIP, _ = ctx.Value(keyRemoteIP).(string)
	return
}
