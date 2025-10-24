package helper

import "context"

type ctxKey string

const (
	UserIDKey ctxKey = "userID"
	AuthKey   ctxKey = "auth"
)

type AuthContext struct {
	UserID    int64
	UserEmail string
	UserName  string
}

func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	v := ctx.Value(AuthKey)
	auth, ok := v.(*AuthContext)
	return auth, ok
}
