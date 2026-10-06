// Package dberr holds the errors the database adapter returns for expected backend conditions.
// It has no internal imports so services can inspect these errors without importing the adapter.
package dberr

import "errors"

// Errors returned for a missing row, a unique violation, and a foreign key violation.
var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("entry already exists")
	ErrForeignKey    = errors.New("foreign key violation")
)
