package pgorm

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type ManagerErrors struct {
}

func (c ManagerErrors) SqlConnections(err error) error {
	if pgErr, ok := err.(*pgconn.ConnectError); ok {
		if strings.Contains(pgErr.Error(), "3D000") {
			// Código de error 3D000 = "invalid_catalog_name" → base de datos no existe
			return fmt.Errorf("base de datos no existe")
		}
		return fmt.Errorf("error PostgreSQL: %s", pgErr.Error())
	}
	// Otro tipo de error
	return fmt.Errorf("error de conexión: %w", err)

}
func (c ManagerErrors) SqlQuery(err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "42P01" {
			return fmt.Errorf("tabla no existe")
		}
		return fmt.Errorf("error PostgreSQL: %s", pgErr.Error())
	}
	// Otro tipo de error
	return fmt.Errorf("error de conexión: %w", err)

}
