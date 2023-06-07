package api

import "net/http"

// Start a transaction
func StartTransaction(rt *_router, w http.ResponseWriter) {
	err := rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the start of the transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Rollback of an already started transaction
func Rollback(rt *_router, w http.ResponseWriter) {
	err := rt.db.Rollback()
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the rollback of the transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Commit of an already started transaction
func Commit(rt *_router, w http.ResponseWriter) {
	err := rt.db.Commit()
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the commit of the transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
