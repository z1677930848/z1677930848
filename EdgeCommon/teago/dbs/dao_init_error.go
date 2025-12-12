package dbs

// daoInitErrorHandlers stores callbacks for DAO init errors.
var daoInitErrorHandlers []func(dao DAOInterface, err error) error

// OnDAOInitError registers a callback for DAO initialization errors.
// In this stub implementation we only store the handler.
func OnDAOInitError(cb func(dao DAOInterface, err error) error) {
	if cb == nil {
		return
	}
	daoInitErrorHandlers = append(daoInitErrorHandlers, cb)
}
