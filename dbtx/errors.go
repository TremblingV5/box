package dbtx

import "errors"

var (
	ErrNotBindWithTx             = errors.New("context not bind with tx")
	ErrNotBindWithTheSpecifiedTx = errors.New("context not bind with the specified tx")
)
