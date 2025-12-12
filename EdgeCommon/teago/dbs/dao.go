package dbs

// DAOObject holds metadata about a DAO.
type DAOObject struct {
	DB       string
	Table    string
	Model    any
	PkName   string
	Instance *DB
}

// DAOWrapper matches the interface expected by NewQuery helper.
type DAOWrapper interface {
	Object() *DAOObject
}

// DAOInterface is kept for backward compatibility.
type DAOInterface interface {
	Object() *DAOObject
}

// DAO is a lightweight DAO wrapper.
type DAO struct {
	DAOObject
}

// NewDAO returns the passed DAO for chaining, matching TeaGo signature.
func NewDAO(dao any) any {
	return dao
}

// Object returns DAOObject.
func (d *DAO) Object() *DAOObject {
	return &d.DAOObject
}

// Query creates a query builder.
func (d *DAO) Query(tx *Tx) *Query {
	return &Query{dao: d, tx: tx}
}

// Query creates a query builder from DAOObject.
func (o *DAOObject) Query(tx *Tx) *Query {
	return &Query{dao: &DAO{DAOObject: *o}, tx: tx}
}

// Object returns self to satisfy DAOWrapper when DAOObject is embedded.
func (o *DAOObject) Object() *DAOObject { return o }

// Exist checks if record exists by primary key (stubbed).
func (o *DAOObject) Exist(tx *Tx, pk any) (bool, error) {
	return o.Query(tx).Pk(pk).Exist()
}

// Save is a placeholder for create/update operations on embedded DAOObject.
func (o *DAOObject) Save(_ *Tx, _ any) error { return nil }

// SaveInt64 is a helper returning int64 id for convenience.
func (o *DAOObject) SaveInt64(_ *Tx, _ any) (int64, error) { return 0, nil }

// Init prepares DAOObject (stubbed).
func (o *DAOObject) Init() error { return nil }

// Delete removes record by primary key (stub).
func (o *DAOObject) Delete(_ *Tx, _ any) (int64, error) { return 0, nil }

// Find retrieves a record by primary key (stub).
func (o *DAOObject) Find(_ *Tx, _ any) (any, error) { return nil, nil }

// Save is a placeholder for create/update operations.
func (d *DAO) Save(tx *Tx, op any) error { return nil }

// Delete deletes record by primary key (stub).
func (d *DAO) Delete(_ *Tx, _ any) (int64, error) { return 0, nil }
