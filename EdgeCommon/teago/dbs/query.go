package dbs

import (
	"github.com/iwind/TeaGo/maps"
)

const (
	QueryJoinLeft  = "LEFT"
	QueryJoinRight = "RIGHT"
)

// Query is a minimal query builder stub to satisfy existing call sites.
type Query struct {
	dao    *DAO
	tx     *Tx
	data   maps.Map
	reuse  bool
	params maps.Map
}

// Attr sets a field equals value.
func (q *Query) Attr(field string, value any) *Query { return q }

// Neq adds not-equal condition.
func (q *Query) Neq(field string, value any) *Query { return q }

// Like adds LIKE condition.
func (q *Query) Like(field string, value any) *Query { return q }

// Param sets a named parameter.
func (q *Query) Param(name string, value any) *Query { return q }

// Set sets a field value for update/insert.
func (q *Query) Set(field string, value any) *Query { return q }

// Sets sets multiple fields at once.
func (q *Query) Sets(data maps.Map) *Query { return q }

// Data sets multiple fields.
func (q *Query) Data(data maps.Map) *Query { return q }

// Where adds where clause (stub).
func (q *Query) Where(cond string, args ...any) *Query { return q }

// In adds IN clause.
func (q *Query) In(field string, values ...any) *Query { return q }

// Table sets table name.
func (q *Query) Table(name string) *Query { return q }

// UseIndex hints using specific indexes.
func (q *Query) UseIndex(_ ...string) *Query { return q }

// Join joins another dao/table.
func (q *Query) Join(_ any, _ string, _ string) *Query { return q }

// Having adds having clause.
func (q *Query) Having(_ string, _ ...any) *Query { return q }

// JSONContains adds JSON containment condition.
func (q *Query) JSONContains(_ string, _ any) *Query { return q }

// Gte adds >= condition.
func (q *Query) Gte(field string, value any) *Query { return q }

// Gt adds > condition.
func (q *Query) Gt(field string, value any) *Query { return q }

// Lte adds <= condition.
func (q *Query) Lte(field string, value any) *Query { return q }

// Lt adds < condition.
func (q *Query) Lt(field string, value any) *Query { return q }

// Between adds BETWEEN condition.
func (q *Query) Between(field string, from any, to any) *Query { return q }

// State adds a custom state flag.
func (q *Query) State(_ ...any) *Query { return q }

// Pk sets primary key value.
func (q *Query) Pk(value any) *Query { return q }

// Reuse marks query reusable.
func (q *Query) Reuse(on bool) *Query {
	q.reuse = on
	return q
}

// Slice sets select fields (stub).
func (q *Query) Slice(target any) *Query { return q }

// Desc orders descending.
func (q *Query) Desc(field string) *Query { return q }

// Asc orders ascending.
func (q *Query) Asc(field string) *Query { return q }

// AscPk orders by primary key ascending.
func (q *Query) AscPk() *Query { return q }

// DescPk orders by primary key descending.
func (q *Query) DescPk() *Query { return q }

// Limit sets limit.
func (q *Query) Limit(count int64) *Query { return q }

// Size sets limit (alias of Limit).
func (q *Query) Size(count int64) *Query { return q }

// Offset sets offset.
func (q *Query) Offset(offset int64) *Query { return q }

// Group sets group by.
func (q *Query) Group(field string) *Query { return q }

// Hint sets sql hint (stub).
func (q *Query) Hint(hint string) *Query { return q }

// Result returns a result helper.
func (q *Query) Result(_ ...any) *Query { return q }

// ResultPk is a shorthand for selecting primary key.
func (q *Query) ResultPk() *Query { return q }

// FindOne returns one row, last id and error.
func (q *Query) FindOne() (maps.Map, int64, error) {
	return nil, 0, nil
}

// Find returns one row (alias).
func (q *Query) Find() (any, error) {
	m, _, err := q.FindOne()
	return m, err
}

// FindStringCol returns a string column with a default.
func (q *Query) FindStringCol(defaultValue string) (string, error) { return defaultValue, nil }

// FindInt64Col returns an int64 column with a default.
func (q *Query) FindInt64Col(defaultValue int64) (int64, error) { return defaultValue, nil }

// FindFloat64Col returns a float64 column with a default.
func (q *Query) FindFloat64Col(defaultValue float64) (float64, error) { return defaultValue, nil }

// FindIntCol returns an int column with a default.
func (q *Query) FindIntCol(defaultValue int) (int, error) { return defaultValue, nil }

// FindBoolCol returns a bool column with a default.
func (q *Query) FindBoolCol(defaultValue ...bool) (bool, error) {
	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}
	return false, nil
}

// FindJSONCol returns JSON bytes with a default.
func (q *Query) FindJSONCol(defaultValue ...[]byte) ([]byte, error) {
	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}
	return nil, nil
}

// FindBytesCol returns raw bytes with a default.
func (q *Query) FindBytesCol(defaultValue ...[]byte) ([]byte, error) {
	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}
	return nil, nil
}

// FindCol returns a column as generic value.
func (q *Query) FindCol(defaultValue ...any) (any, error) {
	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}
	return nil, nil
}

// FindOnes returns multiple rows.
func (q *Query) FindOnes() ([]maps.Map, int64, error) {
	return []maps.Map{}, 0, nil
}

// FindAll wraps FindOnes with total.
func (q *Query) FindAll() ([]any, error) { return []any{}, nil }

// Exist checks existence.
func (q *Query) Exist() (bool, error) { return false, nil }

// Count returns count.
func (q *Query) Count() (int64, error) { return 0, nil }

// CountAttr counts by attribute.
func (q *Query) CountAttr(_ string) (int64, error) { return 0, nil }

// Update updates matched rows.
func (q *Query) Update() (int64, error) { return 0, nil }

// UpdateQuickly updates matched rows (alias).
func (q *Query) UpdateQuickly() error { return nil }

// UpdateOne updates one row.
func (q *Query) UpdateOne() (int64, error) { return 0, nil }

// Insert inserts and returns last id.
func (q *Query) Insert() (int64, error) { return 0, nil }

// InsertIgnore inserts with ignore.
func (q *Query) InsertIgnore() (int64, error) { return 0, nil }

// InsertOrUpdate inserts or updates existing.
func (q *Query) InsertOrUpdate(data maps.Map, updateData ...maps.Map) (int64, int64, error) {
	return 0, 0, nil
}

// InsertOrUpdateQuickly inserts or updates existing.
func (q *Query) InsertOrUpdateQuickly(_ ...maps.Map) error { return nil }

// Replace performs REPLACE INTO.
func (q *Query) Replace(data ...maps.Map) (int64, int64, error) { return 0, 0, nil }

// Delete deletes matched rows.
func (q *Query) Delete() (int64, error) { return 0, nil }

// DeleteQuickly deletes matched rows without returning affected count.
func (q *Query) DeleteQuickly() error { return nil }

// Execute runs raw query.
func (q *Query) Execute() (int64, error) { return 0, nil }

// Sum returns sum of field.
func (q *Query) Sum(field string, _ ...float64) (float64, error) { return 0, nil }

// SumInt64 returns sum as int64.
func (q *Query) SumInt64(field string, _ ...int64) (int64, error) { return 0, nil }

// Max returns max of field.
func (q *Query) Max(field string, _ ...any) (int64, error) { return 0, nil }

// Min returns min of field.
func (q *Query) Min(field string) (any, error) { return nil, nil }

// All returns list of maps.
func (q *Query) All() ([]maps.Map, error) { return []maps.Map{}, nil }

// Avg returns average of field.
func (q *Query) Avg(field string, _ ...float64) (float64, error) { return 0, nil }

// Result helper mirrors TeaGo dbs result set helpers.
type Result struct{}

func (r *Result) FindStringCol(defaultValue string) (string, error) { return defaultValue, nil }
func (r *Result) FindInt64Col(defaultValue int64) (int64, error)    { return defaultValue, nil }
func (r *Result) FindCol(name string) (any, error)                  { return nil, nil }
func (r *Result) FindIntCol(defaultValue int) (int, error)          { return defaultValue, nil }
func (r *Result) FindBoolCol(defaultValue bool) (bool, error)       { return defaultValue, nil }
func (r *Result) Find() (any, error)                                { return nil, nil }
func (r *Result) FindOnes() ([]maps.Map, error)                     { return []maps.Map{}, nil }
func (r *Result) FindAll() ([]any, error)                           { return []any{}, nil }
func (r *Result) Pk(_ ...any) *Result                               { return r }
func (r *Result) Slice(_ any) *Result                               { return r }
func (r *Result) Offset(_ int64) *Result                            { return r }
func (r *Result) Limit(_ int64) *Result                             { return r }
func (r *Result) Attr(_ string, _ any) *Result                      { return r }
func (r *Result) State(_ ...any) *Result                            { return r }
func (r *Result) Where(_ string, _ ...any) *Result                  { return r }
func (r *Result) AscPk() *Result                                    { return r }
func (r *Result) DescPk() *Result                                   { return r }
