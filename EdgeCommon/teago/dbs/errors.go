package dbs

import "errors"

// ErrNotFound indicates no rows found.
var ErrNotFound = errors.New("not found")

// ErrTableNotFound indicates table missing.
var ErrTableNotFound = errors.New("table not found")

// ShowPreparedStatements toggles prepared statement debug output.
var ShowPreparedStatements bool
