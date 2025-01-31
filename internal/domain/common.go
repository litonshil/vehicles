package domain

import (
	"context"
	"vehicles/utils/consts"
)

// ContextWithValue returns a new Context that carries value u.
func ContextWithValue(seedCtx context.Context, key consts.ContextKey, u interface{}) context.Context {
	switch key {
	case consts.ContextKeyUser:
		return context.WithValue(seedCtx, consts.ContextKeyUser, u.(*User))
	}

	return seedCtx
}
