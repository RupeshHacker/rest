package rest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type contextKey string

// UserInfo wraps stringer and uses String to report user name, id and anything else implemented
type UserInfo fmt.Stringer

// GetUserInfo returns user from request context
func GetUserInfo(r *http.Request) (user UserInfo, err error) {

	ctx := r.Context()
	if ctx == nil {
		log.Print("!!!")
		return nil, errors.New("no info about user")
	}
	if u, ok := ctx.Value(contextKey("user")).(UserInfo); ok {
		return u, nil
	}

	return nil, errors.New("user can't be parsed")
}

// SetUserInfo sets user into request context
func SetUserInfo(r *http.Request, user UserInfo) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKey("user"), user)
	return r.WithContext(ctx)
}
