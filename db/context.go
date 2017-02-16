package db

import (
	"database/sql"
	"strings"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/muesli/cct/config"

	log "github.com/Sirupsen/logrus"
	"github.com/muesli/smolder"
)

// APIContext is API's central context
type APIContext struct {
	Config config.Data

	db        *sql.DB
	Queries   []PgQuery
	txIDCount int
}

// APIContextTx is a transactional API conteollyxt
type APIContextTx struct {
	id      int
	context *APIContext
	sqlTx   *sql.Tx
}

// Use sqlAdapter instead of the specific type to allow passing either *sql.DB,
// *sql.Tx, *APIContext or *APIContextTx to a function.
type sqlAdapter interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// PgQuery keeps stats for a single postgres query
type PgQuery struct {
	Query    string
	Duration time.Duration
	TxID     int
}

// NewAPIContext returns a new API context
func (context *APIContext) NewAPIContext() smolder.APIContext {
	ctx := &APIContext{
		db:     GetDatabase(),
		Config: context.Config,
	}
	return ctx
}

// Authentication parses the request for an access-/authtoken and returns the matching user
func (context *APIContext) Authentication(request *restful.Request) (interface{}, error) {
	t := request.QueryParameter("accesstoken")
	if len(t) == 0 {
		t = request.HeaderParameter("authorization")
		if strings.Index(t, " ") > 0 {
			t = strings.TrimSpace(strings.Split(t, " ")[1])
		}
	}

	return context.GetUserByAccessToken(t)
}

// LogSummary logs out the current context stats
func (context *APIContext) LogSummary() {
	for k, v := range context.Queries {
		fields := log.Fields{
			"Num":      k,
			"Query":    v.Query,
			"Duration": v.Duration,
		}
		if v.TxID > -1 {
			fields["Tx"] = v.TxID
		}
		log.WithFields(fields).Debug("Processed postgres query")
	}
}

func (context *APIContext) appendQuery(duration time.Duration, query string, txID int) {
	context.Queries = append(context.Queries, PgQuery{
		Query:    query,
		Duration: duration,
		TxID:     txID,
	})
}

// Exec runs a postgres Exec
func (context *APIContext) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := context.db.Exec(query, args...)

	context.appendQuery(time.Since(start), query, -1)

	return result, err
}

// Query runs a postgres Query
func (context *APIContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := context.db.Query(query, args...)

	context.appendQuery(time.Since(start), query, -1)

	return rows, err
}

// QueryRow runs a postgres QueryRow
func (context *APIContext) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := context.db.QueryRow(query, args...)

	context.appendQuery(time.Since(start), query, -1)

	return row
}

// Begin returns a new API transactional context
func (context *APIContext) Begin() (*APIContextTx, error) {
	start := time.Now()

	tx, err := context.db.Begin()
	if err != nil {
		return nil, err
	}

	txID := context.txIDCount
	context.txIDCount++

	context.appendQuery(time.Since(start), "BEGIN", txID)

	return &APIContextTx{txID, context, tx}, nil
}

// Exec runs a postgres Exec in the transactional context
func (hTx *APIContextTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := hTx.sqlTx.Exec(query, args...)

	hTx.context.appendQuery(time.Since(start), query, hTx.id)

	return result, err
}

// Query runs a postgres Query in the transactional context
func (hTx *APIContextTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := hTx.sqlTx.Query(query, args...)

	hTx.context.appendQuery(time.Since(start), query, hTx.id)

	return rows, err
}

// QueryRow runs a postgres QueryRow in the transactional context
func (hTx *APIContextTx) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := hTx.sqlTx.QueryRow(query, args...)

	hTx.context.appendQuery(time.Since(start), query, hTx.id)

	return row
}

// Commit runs a postgres Commit in the transactional context
func (hTx *APIContextTx) Commit() error {
	start := time.Now()
	err := hTx.sqlTx.Commit()

	hTx.context.appendQuery(time.Since(start), "COMMIT", hTx.id)

	return err
}

// Rollback runs a postgres Rollback in the transactional context
func (hTx *APIContextTx) Rollback() error {
	start := time.Now()
	err := hTx.sqlTx.Rollback()

	hTx.context.appendQuery(time.Since(start), "ROLLBACK", hTx.id)

	return err
}

func (hTx *APIContextTx) commitOrRollbackOnError(err *error) {
	if p := recover(); p != nil {
		hTx.Rollback()
		panic(*err)
	}

	if *err != nil {
		innerErr := hTx.Rollback()
		if innerErr != nil {
			panic(innerErr)
		}
		return
	}

	innerErr := hTx.Commit()
	if innerErr != nil {
		panic(innerErr)
	}
}

/*
func (context *APIContext) Transact(txFunc func(*APIContextTx) error) (err error) {
	tx, err := context.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	return txFunc(tx)
}
*/
