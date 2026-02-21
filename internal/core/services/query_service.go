package services

import (
	"context"

	"github.com/deybin/pgorm/internal/core/builder"
	"github.com/deybin/pgorm/internal/core/clause"
	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/internal/core/ports"
)

type Query struct {
	Schema   string
	Sintaxis *domain.Sintaxis
	Db       ports.DBPort
	Err      error
}

func (q *Query) WorkQueryFull(query string, arg ...interface{}) *Query {
	q.Sintaxis.WorkQueryFull(query, arg...)
	return q
}

func (q *Query) From(table string) *Query {
	q.Sintaxis.From(table)
	return q
}

func (q *Query) Select(campos ...string) *Query {
	q.Sintaxis.Select(campos...)
	return q
}

func (q *Query) Where(where string, op clause.OperatorWhere, arg interface{}) *Query {
	q.Sintaxis.Where(where, op, arg)
	return q
}

func (q *Query) And(and string, op clause.OperatorWhere, arg interface{}) *Query {
	q.Sintaxis.And(and, op, arg)
	return q
}

func (q *Query) Or(or string, op clause.OperatorWhere, arg interface{}) *Query {
	q.Sintaxis.Or(or, op, arg)
	return q
}

func (q *Query) OrderBy(campos ...string) *Query {
	q.Sintaxis.OrderBy(campos...)
	return q
}

func (q *Query) Top(top int) *Query {
	q.Sintaxis.Top(top)
	return q
}

func (q *Query) Limit(limit ...int) *Query {
	q.Sintaxis.Limit(limit...)
	return q
}

func (q *Query) GroupBy(group ...string) *Query {
	q.Sintaxis.GroupBy(group...)
	return q
}

func (q *Query) Join(tp clause.TypeJoin, table string, on string) *Query {
	q.Sintaxis.Join(tp, table, on)
	return q
}

func (q Query) String() string {
	return builder.BuildQuery(q.Sintaxis)
}

func (q *Query) Reset() {
	q.Sintaxis.Reset()
}

/*
Errors devuelve el error almacenado durante la construcción o ejecución de la consulta SQL.

Esta función permite recuperar el error que se haya producido en alguna etapa del proceso de construcción,
ejecución o procesamiento de la consulta SQL. Es útil para el manejo de errores encadenado, ya que muchas funciones
devuelven el mismo struct `Query` y almacenan el error internamente.

Ejemplo de uso:

	queryBuilder :=  new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetTable("my_table").Select("campo1").Where("campo2", "=", valor).Exec()
	if err := queryBuilder.Errors(); err != nil {
	    log.Println("Ocurrió un error:", err)
	}

Devuelve:
  - El error almacenado, si existe; de lo contrario, devuelve `nil`.
*/
func (q *Query) Errors() error {
	return q.Err
}

/*
Close libera los recursos asociados con la consulta SQL.

Esta función cierra el conjunto de resultados (`rowSql`) y la conexión (`conn`) asociada al struct Query.
Debe ser llamada después de ejecutar una consulta para liberar correctamente los recursos y evitar fugas de memoria o conexiones abiertas.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetTable("my_table").Select("campo1").Exec()
	defer queryBuilder.Close()

Importante:
  - Si `rowSql` o `conn` no fueron inicializados o ya fueron cerrados, esta función podría generar un pánico.
    Se recomienda validar o asegurar su existencia antes de usar esta función si se usa fuera del flujo estándar.

No devuelve ningún valor.
*/
func (q *Query) Close() {
	// q.conn.Close()
}

func (q *Query) Exec() ([]map[string]any, error) {
	result, err := q.Db.Execute(q.Schema, q.Db.Context(), q.String(), q.Sintaxis.Arguments()...)
	if err != nil {
		return result, err
	}
	return result, err
}

func (q *Query) Procedure() error {
	err := q.Db.Procedure(q.Schema, q.Db.Context(), q.String(), q.Sintaxis.Arguments()...)
	return err
}

func (q *Query) ExecWithContext(ctx context.Context) ([]map[string]any, error) {
	result, err := q.Db.Execute(q.Schema, ctx, q.String(), q.Sintaxis.Arguments()...)
	if err != nil {
		return result, err
	}
	return result, err
}

func (q *Query) ProcedureWithContext(ctx context.Context) error {
	err := q.Db.Procedure(q.Schema, ctx, q.String(), q.Sintaxis.Arguments()...)
	return err
}

func QueryExec[T any](q *Query) (T, error) {
	var dest T
	err := q.Db.ExecuteWithPgxScan(q.Schema, q.Db.Context(), &dest, q.String(), q.Sintaxis.Arguments()...)
	return dest, err
}

func QueryExecWithContext[T any](ctx context.Context, q *Query) (T, error) {
	var dest T
	err := q.Db.ExecuteWithPgxScan(q.Schema, ctx, &dest, q.String(), q.Sintaxis.Arguments()...)
	return dest, err
}
