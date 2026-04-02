package ports

import (
	"context"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPort interface {
	Execute(ctx context.Context, sql string, args ...any) ([]map[string]any, error)
	ExecuteWithPgxScan(ctx context.Context, dest any, sql string, args ...any) error
	ExecuteWithPgxScanAndSchema(schema string, ctx context.Context, dest any, sql string, args ...any) error
	Procedure(ctx context.Context, sql string, args ...any) error
	ProcedureWithSchema(schema string, ctx context.Context, sql string, arguments ...any) error
	ExecuteTransactions(ctx context.Context, dataExec ...adapters.DataExec) error
	ExecuteTransactionsWithSchema(schema string, ctx context.Context, dataExec ...adapters.DataExec) error
	ExecuteTransactionsMulti(ctx context.Context, dataExec ...[]adapters.DataExec) error
	ExecuteTransactionsMultiWithSchema(schema string, ctx context.Context, dataExec ...[]adapters.DataExec) error
	Pool() *pgxpool.Pool
}
