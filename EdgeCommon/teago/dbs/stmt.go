package dbs

// Stmt is a prepared statement stub.
type Stmt struct {
	db    *DB
	query string
}

// Close closes statement.
func (s *Stmt) Close() error { return nil }

// FindCol executes statement and returns a column value (stubbed).
func (s *Stmt) FindCol(index int, args ...any) (any, error) {
	if s != nil && s.db != nil {
		return s.db.FindCol(index, s.query, args...)
	}
	return nil, nil
}
