package adapters

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

type Actions uint8

const (
	NONE Actions = iota
	INSERT
	UPDATE
	DELETE
)

type DataExec struct {
	Querys string
	Values []any
	Action Actions
}

// dbExecutor es una interfaz interna para aceptar tanto conexiones como transacciones
type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}
