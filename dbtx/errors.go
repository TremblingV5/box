package dbtx

import "errors"

var (
	// ErrNotBindWithTx is returned when trying to get a transaction from a context
	// but no transaction has been bound to that context
	ErrNotBindWithTx = errors.New("context not bind with tx")

	// ErrNotBindWithTheSpecifiedTx is returned when trying to get a transaction of a specific type
	// but the transaction bound to the context is not of that type
	ErrNotBindWithTheSpecifiedTx = errors.New("context not bind with the specified tx")
)
