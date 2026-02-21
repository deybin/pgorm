package pgorm

import (
	"context"

	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/internal/core/ports"
	"github.com/deybin/pgorm/internal/core/services"
	"github.com/deybin/pgorm/migrator"
)

// QUERY
func NewQuery(db ports.DBPort) *services.Query {
	return &services.Query{
		Sintaxis: domain.NewSintaxis(), // Asumiendo que tienes un constructor en domain
		Db:       db,
	}
}

func NewQueryWithSchema(schema string, db ports.DBPort) *services.Query {
	return &services.Query{
		Schema:   schema,
		Sintaxis: domain.NewSintaxis(), // Asumiendo que tienes un constructor en domain
		Db:       db,
	}
}

func QueryExec[T any](q *services.Query) (T, error) {
	var dest T
	err := q.Db.ExecuteWithPgxScan(q.Schema, q.Db.Context(), &dest, q.String(), q.Sintaxis.Arguments()...)
	q.Sintaxis.Reset()
	return dest, err
}

func QueryExecWithContext[T any](ctx context.Context, q *services.Query) (T, error) {
	var dest T
	err := q.Db.ExecuteWithPgxScan(q.Schema, ctx, &dest, q.String(), q.Sintaxis.Arguments()...)
	q.Sintaxis.Reset()
	return dest, err
}

// CRUD
func NewSqlExecSingles(db ports.DBPort, s migrator.Schema, datos ...migrator.Entity) *services.SqlExecSingles {
	return &services.SqlExecSingles{
		Transactions: domain.NewTransaction(s, datos...), // Asumiendo que tienes un constructor en domain
		Db:           db,
	}
}

func NewSqlExecSinglesWithSchema(schema string, db ports.DBPort, s migrator.Schema, datos ...migrator.Entity) *services.SqlExecSingles {
	return &services.SqlExecSingles{
		Transactions: domain.NewTransaction(s, datos...), // Asumiendo que tienes un constructor en domain
		Db:           db,
		Schemas:      schema,
	}
}

func TransactionExec(s *services.SqlExecSingles) error {
	return s.Db.ExecuteTransactions(s.Schemas, s.Db.Context(), s.Transactions.Query()...)
}

func TransactionExecWithContext(ctx context.Context, s *services.SqlExecSingles) error {
	return s.Db.ExecuteTransactions(s.Schemas, ctx, s.Transactions.Query()...)
}

func NewSqlExecMultiples(db ports.DBPort) *services.SqlExecMultiples {
	return &services.SqlExecMultiples{
		Db: db,
	}
}

func NewSqlExecMultiplesWithSchema(db ports.DBPort, schema string) *services.SqlExecMultiples {
	return &services.SqlExecMultiples{
		Db:      db,
		Schemas: schema,
	}
}

func TransactionMultiExec(s *services.SqlExecMultiples) error {
	return s.Db.ExecuteTransactionsWithContext(s.Schemas, s.Db.Context(), s.DataExec()...)
}

func TransactionMultiExecWithContext(ctx context.Context, s *services.SqlExecMultiples) error {
	return s.Db.ExecuteTransactionsWithContext(s.Schemas, ctx, s.DataExec()...)
}

//Generator

func NewEntitySchema(database string, table string) {
	migrator.GenerateSchemaFile(database, table)
}
