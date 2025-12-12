package dbs

// Table describes a database table schema.
// This is a lightweight stub used by setup/sql_dump.
type Table struct {
	Name      string
	Engine    string
	Collation string
	Code      string
	Fields    []*TableField
	Indexes   []*TableIndex
}

// TableField describes a table column.
type TableField struct {
	Name string
	Code string
}

// Definition returns field definition SQL.
func (f *TableField) Definition() string { return f.Code }

// TableIndex describes an index.
type TableIndex struct {
	Name string
	Code string
}

// Definition returns index definition SQL.
func (i *TableIndex) Definition() string { return i.Code }
