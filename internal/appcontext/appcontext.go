// Package auth provides functions for working with application context.
package appcontext

import "context"

type key int

const userIDKey key = 0

//ContextWithUserID adds user id to application context.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

//UserIDFromContext gets user id from application context.
func UserIDFromContext(ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}
