package pgorm

import (
	"context"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/internal/core/ports"
	"github.com/deybin/pgorm/internal/core/services"
	"github.com/deybin/pgorm/migrator"
)

//Conexión

func NewPool() (*adapters.PgxAdapter, error) {
	return adapters.NewPool(adapters.ConfigPgxAdapter{})
}

func NewPoolWithConfig(config adapters.ConfigPgxAdapter) (*adapters.PgxAdapter, error) {
	return adapters.NewPoolWithConfig(config)
}

// QUERY
func NewQuery() *services.Query {
	return &services.Query{
		Sintaxis: domain.NewSintaxis(),
	}
}

func ExecQuery[T any](db ports.DBPort, ctx context.Context, q *services.Query) (T, error) {
	var dest T
	err := db.ExecuteWithPgxScan(ctx, &dest, q.String(), q.Sintaxis.Arguments()...)
	q.Sintaxis = &domain.Sintaxis{}
	return dest, err
}

func ExecQueryWithSchema[T any](db ports.DBPort, schema string, ctx context.Context, q *services.Query) (T, error) {
	var dest T
	err := db.ExecuteWithPgxScanAndSchema(schema, ctx, &dest, q.String(), q.Sintaxis.Arguments()...)
	q.Sintaxis.Reset()
	return dest, err
}

//Procedure

func ExecProcedure(db ports.DBPort, ctx context.Context, q *services.Query) error {
	err := db.Procedure(ctx, q.String(), q.Sintaxis.Arguments()...)
	return err
}

func ExecProcedureWithSchema(db ports.DBPort, schema string, ctx context.Context, q *services.Query) error {
	err := db.ProcedureWithSchema(schema, ctx, q.String(), q.Sintaxis.Arguments()...)
	return err
}

// CRUD SINGLE
func NewSqlExecSingles(s migrator.Schema, datos ...migrator.Entity) *services.SqlExecSingles {
	return &services.SqlExecSingles{
		Transactions: domain.NewTransaction(s, datos...),
	}
}

func ExecTransaction(db ports.DBPort, ctx context.Context, s *services.SqlExecSingles) error {
	return db.ExecuteTransactions(ctx, s.Transactions.Query()...)
}

func ExecTransactionWithSchema(db ports.DBPort, schema string, ctx context.Context, s *services.SqlExecSingles) error {
	return db.ExecuteTransactionsWithSchema(schema, ctx, s.Transactions.Query()...)
}

// CRUD MULTI

func NewSqlExecMultiples() *services.SqlExecMultiples {
	return &services.SqlExecMultiples{}
}

func ExecTransactionMulti(db ports.DBPort, ctx context.Context, s *services.SqlExecMultiples) error {
	return db.ExecuteTransactionsMulti(ctx, s.DataExec()...)
}

func ExecTransactionMultiWithSchema(db ports.DBPort, schema string, ctx context.Context, s *services.SqlExecMultiples) error {
	return db.ExecuteTransactionsMultiWithSchema(schema, ctx, s.DataExec()...)
}

//Generator

func NewEntitySchema(database string, table string) {
	migrator.GenerateSchemaFile(database, table)
}

func NewEntitySchemaWithSchemaName(database string, table string, schema string) {
	migrator.GenerateSchemaFileWithSchema(database, table, schema)
}
