package pgorm

import (
	"log"

	"github.com/deybin/pgorm/clause"
	"github.com/deybin/pgorm/internal"
	"github.com/deybin/pgorm/logger"
	"github.com/deybin/pgorm/schema"
	"github.com/jackc/pgx/v5"
)

type QueryG[T any] struct {
	query  Query
	model  T
	result []T
}

func G[T any](config Config) *QueryG[T] {
	var err error
	var q = QueryG[T]{}
	// q.model = T
	if config.Cloud {
		if q.query.conn, err = new(internal.Connection).New("").PoolMaster(); err != nil {
			q.query.err = err
			return &q
		}
	} else {
		if q.query.conn, err = new(internal.Connection).New(config.Database).NewPool(); err != nil {
			q.query.err = err
			return &q
		}
	}

	return &q
}

func (q *QueryG[T]) WorkQueryFull(query string, arg ...interface{}) *QueryG[T] {
	q.query.WorkQueryFull(query, arg...)
	return q
}

func (q *QueryG[T]) From(table schema.Schema) *QueryG[T] {
	q.query.From(table)
	return q
}

func (q *QueryG[T]) Select(campos ...string) *QueryG[T] {
	q.query.Select(campos...)
	return q
}

func (q *QueryG[T]) Where(where string, op clause.OperatorWhere, arg interface{}) *QueryG[T] {
	q.query.Where(where, op, arg)
	return q
}

func (q *QueryG[T]) And(and string, op clause.OperatorWhere, arg interface{}) *QueryG[T] {
	q.query.And(and, op, arg)
	return q
}

func (q *QueryG[T]) Or(or string, op clause.OperatorWhere, arg interface{}) *QueryG[T] {
	q.query.Or(or, op, arg)
	return q
}

func (q *QueryG[T]) OrderBy(campos ...string) *QueryG[T] {
	q.query.OrderBy(campos...)
	return q
}

func (q *QueryG[T]) Top(top int) *QueryG[T] {
	q.query.Top(top)
	return q
}

func (q *QueryG[T]) Limit(limit ...int) *QueryG[T] {
	q.query.Limit(limit...)
	return q
}

func (q *QueryG[T]) GroupBy(group ...string) *QueryG[T] {
	q.query.GroupBy(group...)
	return q
}

func (q *QueryG[T]) Join(tp clause.TypeJoin, table schema.Schema, on string) *QueryG[T] {
	q.query.Join(tp, table, on)
	return q
}

func (q *QueryG[T]) Exec() *QueryG[T] {
	if q.query.err != nil {
		return q
	}

	q.query.sessionActiva = false
	queryString := q.query.build()
	// fmt.Println("query:", queryString, q.query.query.args)

	rows, err := q.query.conn.Pool().Query(q.query.conn.Context(), queryString, q.query.query.args...)
	if err != nil {
		q.query.err = err
		return q
	}

	q.query.rowSql = rows

	return q
}

func (q *QueryG[T]) ExecCtx() *QueryG[T] {
	if q.query.err != nil {
		return q
	}

	q.query.sessionActiva = true
	queryString := q.query.build()
	// fmt.Println("query:", queryString)

	rows, err := q.query.conn.Pool().Query(q.query.conn.Context(), queryString, q.query.query.args...)
	if err != nil {
		q.query.err = err
		return q
	}

	q.query.rowSql = rows

	return q
}

func (q *QueryG[T]) Procedure() error {
	return q.query.Procedure()
}

func (q *QueryG[T]) ProcedureCtx() error {
	return q.query.ProcedureCtx()
}

func (q *QueryG[T]) One() (T, error) {

	err := q.obtenerDatos()

	if err != nil {
		log.Printf("check failed: %v", err)
		return q.model, logger.ManagerErrors{}.SqlQuery(err)
	}

	return q.result[0], nil

}

func (q *QueryG[T]) Value(extractor func(T) any) (any, error) {
	err := q.obtenerDatos()
	if err != nil {

		log.Printf("check failed: %v", err)
		return nil, logger.ManagerErrors{}.SqlQuery(err)
	}

	return extractor(q.result[0]), nil
}

func (q *QueryG[T]) All() ([]T, error) {
	err := q.obtenerDatos()
	if err != nil {
		log.Printf("check failed: %v", err)
		return []T{}, logger.ManagerErrors{}.SqlQuery(err)
	}
	return q.result, nil
}

func (q *QueryG[T]) String() string {
	return q.query.String()
}

func (q *QueryG[T]) Errors() error {
	return q.query.Errors()
}

func (q *QueryG[T]) Close() {
	q.query.Close()
}

func (q *QueryG[T]) Reset() {
	q.query.Reset()
	q.result = []T{}
}

func (q *QueryG[T]) obtenerDatos() error {
	if q.query.err != nil {
		log.Printf("check failed: %v", q.query.err)
		return logger.ManagerErrors{}.SqlQuery(q.query.err)
	}
	if !q.query.sessionActiva {
		defer q.query.conn.Close()
	}
	defer q.query.rowSql.Close()
	var err error
	q.result, err = pgx.CollectRows(q.query.rowSql, pgx.RowToStructByName[T])

	if err != nil {
		log.Printf("check failed: %v", err)
		return logger.ManagerErrors{}.SqlQuery(err)
	}
	return nil

}
