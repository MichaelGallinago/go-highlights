package grpc

import "errors"

var (
	ErrNotFound = errors.New("item not found")
	ErrInternal = errors.New("internal error")
)
