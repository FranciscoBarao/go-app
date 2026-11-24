// Package store holds the errors a database adapter returns for expected backend conditions.
// It has no internal imports, so a service can inspect these errors without importing the adapter.
package store

import "errors"

// Errors returned for a missing row, a unique violation, and a foreign key violation.
var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("entry already exists")
	ErrForeignKey    = errors.New("foreign key violation")
)
