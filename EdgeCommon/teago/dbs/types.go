package dbs

// JSON represents JSON blob content.
type JSON []byte

// FieldName is an alias for field identifiers.
type FieldName = string

// SQL wraps raw SQL expressions.
type SQL string

// IsNull reports whether JSON content is empty or "null".
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// IsNotNull reports whether JSON content is not empty or null.
func (j JSON) IsNotNull() bool { return !j.IsNull() }
