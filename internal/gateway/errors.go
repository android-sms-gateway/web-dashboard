package gateway

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUpstream     = errors.New("upstream error")
)
