package ports

import (
	"context"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPort interface {
	Execute(schema string, ctx context.Context, sql string, args ...any) ([]map[string]any, error)
	ExecuteWithPgxScan(schema string, ctx context.Context, dest any, sql string, args ...any) error
	Procedure(schema string, ctx context.Context, sql string, args ...any) error
	ExecuteTransactions(schema string, ctx context.Context, dataExec ...adapters.DataExec) error
	ExecuteTransactionsWithContext(schema string, ctx context.Context, dataExec ...[]adapters.DataExec) error
	Context() context.Context
	Pool() *pgxpool.Pool
}
