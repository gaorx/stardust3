package sdcontext

import (
	"context"
	"time"
)

func OrBackground(ctxOpt context.Context) context.Context {
	if ctxOpt == nil {
		return context.Background()
	} else {
		return ctxOpt
	}
}

func OrTodo(ctxOpt context.Context) context.Context {
	if ctxOpt == nil {
		return context.TODO()
	} else {
		return ctxOpt
	}
}

func OrTimeout(ctxOpt context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(OrBackground(ctxOpt), timeout)
}

func OrDeadline(ctxOpt context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(OrBackground(ctxOpt), d)
}
