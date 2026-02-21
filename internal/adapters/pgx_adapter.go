package adapters

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/deybin/pgorm/internal/utils"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type PgxAdapter struct {
	db       *pgxpool.Pool
	context  context.Context
	database string
	schema   string
	server   string
	user     string
	password string
	port     string
	ssl      string //disable,require
	appName  string
}

type ConfigPgxAdapter struct {
	Host            string
	Database        string
	Schema          string
	User            string
	Password        string
	Port            string
	Ssl             string //disable,require
	AppName         string
	MaxConns        int32
	MinConns        int32
	MaxConnIdleTime time.Duration
}

func NewPool(setting ConfigPgxAdapter) (*PgxAdapter, error) {
	_ = godotenv.Load()

	var dsn strings.Builder

	dsn.WriteString("postgres://")

	user := os.Getenv("ENV_DB_USER")

	if user != "" {
		dsn.WriteString(user)
		dsn.WriteByte(':')
	}

	password := os.Getenv("ENV_DB_PASSWORD")

	if password != "" {
		dsn.WriteString(password)
	}

	dsn.WriteByte('@')

	host := os.Getenv("ENV_DB_SERVER")

	if host != "" {
		dsn.WriteString(host)
		dsn.WriteByte(':')
	}

	port := os.Getenv("ENV_DB_PORT")

	if port != "" {
		dsn.WriteString(port)
	}

	dsn.WriteByte('/')

	database := os.Getenv("ENV_DB_DATABASE")

	if database != "" {
		dsn.WriteString(database)
	}

	var arg []string

	schema := os.Getenv("ENV_DB_SCHEMA")

	if schema != "" {
		arg = append(arg, fmt.Sprintf("search_path=%s", schema))
	}

	ssl := os.Getenv("ENV_DB_SSL")

	if ssl != "" {
		arg = append(arg, fmt.Sprintf("sslmode=%s", ssl))
	}

	appName := os.Getenv("ENV_DB_APP")

	if appName != "" {
		arg = append(arg, fmt.Sprintf("application_name=%s", appName))
	}

	if len(arg) > 0 {
		dsn.WriteByte('?')
		dsn.WriteString(strings.Join(arg, "&"))
	}

	config, _ := pgxpool.ParseConfig(dsn.String())
	useConfig := false

	if !reflect.ValueOf(setting.MaxConns).IsZero() {
		config.MaxConns = setting.MaxConns // Máximo 50 conexiones físicas
		useConfig = true
	}
	if !reflect.ValueOf(setting.MinConns).IsZero() {
		config.MinConns = setting.MinConns // Mantener 5 siempre listas
		useConfig = true
	}
	if !reflect.ValueOf(setting.MaxConnIdleTime).IsZero() {
		config.MaxConnIdleTime = setting.MaxConnIdleTime // Cerrar conexiones que nadie usa
		useConfig = true
	}
	ctx := context.Background()

	var pool *pgxpool.Pool
	var err error
	if useConfig {
		pool, err = pgxpool.NewWithConfig(ctx, config)
	} else {
		pool, err = pgxpool.New(ctx, dsn.String())
	}

	if err != nil {
		slog.Error("Fallo en la conexión", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Fallo ping", "error", err)
		return nil, err
	}

	return &PgxAdapter{
		db:       pool,
		context:  ctx,
		database: database,
		schema:   schema,
		server:   host,
		user:     user,
		password: password,
		ssl:      ssl,
		port:     port,
		appName:  appName,
	}, nil

}

func NewPoolWithConfig(setting ConfigPgxAdapter) (*PgxAdapter, error) {
	_ = godotenv.Load()

	var dsn strings.Builder

	dsn.WriteString("postgres://")

	user := os.Getenv("ENV_DB_USER")
	if !reflect.ValueOf(setting.User).IsZero() {
		user = setting.User
	}
	if user != "" {
		dsn.WriteString(user)
		dsn.WriteByte(':')
	}
	password := os.Getenv("ENV_DB_PASSWORD")
	if !reflect.ValueOf(setting.Password).IsZero() {
		password = setting.Password
	}
	if password != "" {
		dsn.WriteString(password)
	}

	dsn.WriteByte('@')

	host := os.Getenv("ENV_DB_SERVER")
	if !reflect.ValueOf(setting.Host).IsZero() {
		host = setting.Host
	}
	if host != "" {
		dsn.WriteString(host)
		dsn.WriteByte(':')
	}

	port := os.Getenv("ENV_DB_PORT")
	if !reflect.ValueOf(setting.Port).IsZero() {
		port = setting.Port
	}
	if port != "" {
		dsn.WriteString(port)
	}

	dsn.WriteByte('/')

	database := os.Getenv("ENV_DB_DATABASE")
	if !reflect.ValueOf(setting.Database).IsZero() {
		database = setting.Database
	}
	if database != "" {
		dsn.WriteString(database)
	}

	var arg []string

	schema := os.Getenv("ENV_DB_SCHEMA")
	if !reflect.ValueOf(setting.Schema).IsZero() {
		schema = setting.Schema
	}
	if schema != "" {
		arg = append(arg, fmt.Sprintf("search_path=%s", schema))
	}

	ssl := os.Getenv("ENV_DB_SSL")
	if !reflect.ValueOf(setting.Ssl).IsZero() {
		ssl = setting.Ssl
	}
	if ssl != "" {
		arg = append(arg, fmt.Sprintf("sslmode=%s", ssl))
	}

	appName := os.Getenv("ENV_DB_APP")
	if !reflect.ValueOf(setting.AppName).IsZero() {
		appName = setting.AppName
	}
	if appName != "" {
		arg = append(arg, fmt.Sprintf("application_name=%s", appName))
	}

	if len(arg) > 0 {
		dsn.WriteByte('?')
		dsn.WriteString(strings.Join(arg, "&"))
	}

	config, _ := pgxpool.ParseConfig(dsn.String())
	useConfig := false

	if !reflect.ValueOf(setting.MaxConns).IsZero() {
		config.MaxConns = setting.MaxConns // Máximo 50 conexiones físicas
		useConfig = true
	}
	if !reflect.ValueOf(setting.MinConns).IsZero() {
		config.MinConns = setting.MinConns // Mantener 5 siempre listas
		useConfig = true
	}
	if !reflect.ValueOf(setting.MaxConnIdleTime).IsZero() {
		config.MaxConnIdleTime = setting.MaxConnIdleTime // Cerrar conexiones que nadie usa
		useConfig = true
	}
	ctx := context.Background()

	var pool *pgxpool.Pool
	var err error
	if useConfig {
		pool, err = pgxpool.NewWithConfig(ctx, config)
	} else {
		pool, err = pgxpool.New(ctx, dsn.String())
	}

	if err != nil {
		slog.Error("Fallo en la conexión", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Fallo ping", "error", err)
		return nil, err
	}

	return &PgxAdapter{
		db:       pool,
		context:  ctx,
		database: database,
		schema:   schema,
		server:   host,
		user:     user,
		password: password,
		ssl:      ssl,
		port:     port,
		appName:  appName,
	}, nil

}

func (c PgxAdapter) Pool() *pgxpool.Pool {
	return c.db
}

func (c PgxAdapter) Context() context.Context {
	return c.context
}

func (c *PgxAdapter) Close() {
	c.db.Close()
}

func (p PgxAdapter) setSchema(ctx context.Context, conn *pgxpool.Conn, schema string) error {
	if schema != "" {
		if _, err := conn.Exec(ctx, "SET search_path TO "+schema); err != nil {
			slog.Error("Fallo al acceder al schema", "error", err)
			return err
		}
	}
	return nil
}

func (p PgxAdapter) executeInternal(ctx context.Context, exec dbExecutor, data ...DataExec) error {
	cross := false // O tu lógica de negocio para cross
	for _, item := range data {
		sqlPre := item.Querys
		if cross && item.Action == UPDATE {
			sqlPre = utils.QueryCrossUpdate(sqlPre)
		}

		if _, err := exec.Exec(ctx, sqlPre, item.Values...); err != nil {
			slog.Error("Fallo Exec", "error", err)
			return err
		}
	}
	return nil
}

/*
keyFieldName extrae los nombres de las columnas devueltos por la consulta ejecutada y los almacena en `q.colSql`.

Esta función utiliza los metadatos de la fila (`pgx.Rows.FieldDescriptions()`) para recuperar el nombre de cada columna del
resultado actual, lo cual es útil para operaciones dinámicas como el mapeo a `map[string]interface{}` en `builderResult()`.

Debe llamarse después de ejecutar una consulta `SELECT`, y antes de intentar leer los resultados si se planea acceder a
las columnas por nombre.

Ejemplo de uso:

	query.ExecCtx()
	query.keyFieldName()

Retorna:
  - No retorna valor, pero actualiza internamente `q.colSql` con los nombres de columnas actuales.
*/
func (q *PgxAdapter) keyFieldName(rows pgx.Rows) []string {
	fieldDescription := rows.FieldDescriptions()
	colSql := make([]string, len(fieldDescription))
	for i, fd := range fieldDescription {
		colSql[i] = string(fd.Name)
	}
	return colSql
}

/*
Exec ejecuta la consulta SQL construida utilizando pgx y almacena los resultados en la estructura Query.

Esta función se encarga de ejecutar la consulta previamente construida con los métodos del builder (Select, Where, etc.),
o mediante una consulta SQL personalizada con `SetQueryString`. Soporta consultas normales (`SELECT`) y procedimientos (`procedure`).

Si es una consulta `SELECT`, guarda los resultados (`pgx.Rows`) y los nombres de las columnas para su posterior lectura.
Si es una ejecución (`INSERT`, `UPDATE`, `DELETE`, etc.) sin retorno de resultados, simplemente la ejecuta.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetTable("my_table").Select("id, name").Where("status", "=", "active").Exec()

Devuelve:
  - Un puntero al struct Query actualizado, incluyendo los resultados o el error si ocurrió alguno.
*/

func (p PgxAdapter) Execute(schema string, ctx context.Context, sql string, args ...any) ([]map[string]any, error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		slog.Error("Fallo conexión db", "error", err)
		return []map[string]any{}, err
	}
	defer conn.Release() // SEGURO: Siempre vuelve al pool al terminar la función

	// 4. Configurar el esquema para ESTA sesión
	if schema != "" {
		_, err = conn.Exec(ctx, "SET search_path TO "+schema)
		if err != nil {
			slog.Error("Fallo al acceder al schema", "error", err)
			return []map[string]any{}, err
		}
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return []map[string]any{}, err
	}

	defer rows.Close()

	cols := p.keyFieldName(rows)

	fieldDescs := rows.FieldDescriptions()
	result := make([]map[string]any, 0)
	for rows.Next() {
		row, err := p.builderResult(cols, rows)
		if err != nil {
			return []map[string]any{}, err
		}
		row = p.normalizeRow(row, fieldDescs)
		result = append(result, row)
	}

	return result, nil
}

func (p PgxAdapter) ExecuteWithPgxScan(schema string, ctx context.Context, dest any, sql string, args ...any) error {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		slog.Error("Fallo conexión db", "error", err)
		return err
	}
	defer conn.Release() // SEGURO: Siempre vuelve al pool al terminar la función

	// 4. Configurar el esquema para ESTA sesión
	if schema != "" {
		_, err = conn.Exec(ctx, "SET search_path TO "+schema)
		if err != nil {
			slog.Error("Fallo al acceder al schema", "error", err)
			return err
		}
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rv := reflect.ValueOf(dest)
	// 1. Primero verificamos si es un puntero (porque dest siempre debe ser &variable)
	if rv.Kind() == reflect.Ptr {
		// 2. Obtenemos el elemento al que apunta el puntero
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		return pgxscan.ScanAll(dest, rows)

	}
	// Para un solo elemento, evitamos ScanOne para que no explote si hay > 1
	if rows.Next() {
		return pgxscan.NewRowScanner(rows).Scan(dest)
	}

	return fmt.Errorf("not found information")
}

/*
Procedure ejecuta una consulta de tipo procedimiento (por ejemplo, `CALL` o funciones sin retorno de datos) utilizando `pgx`.

Esta función está diseñada para ejecutar procedimientos almacenados u operaciones SQL que **no devuelven resultados**
(solo efectos colaterales en la base de datos).

Comportamiento:
- Si existe un error anterior (`q.err`), lo retorna inmediatamente sin ejecutar la consulta.
- Desactiva el indicador `sessionActiva`, indicando que no se espera continuar usando la misma conexión.
- Ejecuta la consulta generada con `getQuery()` y los argumentos acumulados.
- Si ocurre un error al ejecutar, se almacena y se retorna.
- Siempre cierra la conexión (`q.conn.Close()`) al finalizar.

Ejemplo de uso:

	query := new(pgorm.Query).New(pgorm.QConfig{Database: "my_db"})
	query.SetQueryString("CALL procesar_datos($1)").SetArgs(123).Procedure()

Retorna:
  - `nil` si la ejecución fue exitosa,
  - Un `error` si la ejecución falló o si ya existía un error previo en el estado del `Query`.
*/
func (p *PgxAdapter) Procedure(schema string, ctx context.Context, sql string, arguments ...any) error {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		slog.Error("Fallo conexión db", "error", err)
		return err
	}
	defer conn.Release() // SEGURO: Siempre vuelve al pool al terminar la función

	// 4. Configurar el esquema para ESTA sesión
	if schema != "" {
		_, err = conn.Exec(ctx, "SET search_path TO "+schema)
		if err != nil {
			slog.Error("Fallo al acceder al schema", "error", err)
			return err
		}
	}
	if _, err := conn.Exec(ctx, sql, arguments...); err != nil {
		slog.Error("Fallo Exec", "error", err)
		return err
	}
	return nil

}

/*
builderResult extrae una fila del resultado actual (`q.rowSql`) y la convierte en un mapa con claves como los nombres de las columnas
y valores asociados a cada campo.

Este método está pensado para usarse después de ejecutar una consulta `SELECT` mediante `Exec()` o `ExecCtx()`, y asume que
`q.colSql` contiene los nombres de las columnas y `q.rowSql` contiene un cursor (`pgx.Rows`) activo sobre los resultados.

Funcionamiento:
- Crea slices de interfaces para capturar los datos de cada columna.
- Usa `Scan` para copiar los valores de la fila actual del cursor en esos punteros.
- Construye un `map[string]interface{}` usando los nombres de columna como claves y los valores extraídos como valores.

Si ocurre un error durante el `Scan`, lo registra en los logs y devuelve un error personalizado con `ManagerErrors{}.SqlQuery(err)`.

Ejemplo de uso:

	query.ExecCtx()
	for query.rowSql.Next() {
		result, err := query.builderResult()
		if err != nil {
			// manejar error
		}
		fmt.Println(result["id"], result["nombre"])
	}

Retorna:
  - Un `map[string]interface{}` que representa una fila del resultado, donde cada clave es el nombre de la columna.
  - Un `error` si falla el escaneo de los datos.
*/
func (p PgxAdapter) builderResult(cols []string, row pgx.Rows) (map[string]interface{}, error) {

	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := row.Scan(columnPointers...); err != nil {
		slog.Error("Fallo row scan", "error", err)
		return map[string]interface{}{}, err
	}

	//Crea nuestro mapa y recupera el valor de cada columna del segmento de punteros, almacenándolo en el mapa con el nombre de la columna como clave.
	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		l := *val
		if l != nil {

			m[colName] = l

		} else {
			m[colName] = l
		}
	}

	return m, nil

}

// normalizeRow convierte UUIDs (OID 2950) a string
func (p PgxAdapter) normalizeRow(row map[string]interface{}, fieldDescs []pgconn.FieldDescription) map[string]interface{} {
	for _, fd := range fieldDescs {
		colName := string(fd.Name)

		val, exists := row[colName]
		if !exists {
			continue
		}

		// Detectamos si el tipo es UUID (OID = 2950)
		if fd.DataTypeOID == 2950 {

			var b []byte

			switch v := val.(type) {
			case [16]byte:
				b = v[:] // array → slice
			case []byte:
				b = v // slice directo
			default:
				continue // no es UUID válido
			}

			if u, err := uuid.FromBytes(b); err == nil {
				row[colName] = u.String()
			}

		}
	}
	return row
}

func (p PgxAdapter) ExecuteTransactions(schema string, ctx context.Context, dataExec ...DataExec) error {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		slog.Error("Fallo conexión db", "error", err)
		return err
	}
	defer conn.Release() // SEGURO: Siempre vuelve al pool al terminar la función

	// 4. Configurar el esquema para ESTA sesión
	if err := p.setSchema(ctx, conn, schema); err != nil {
		return err
	}

	return p.executeInternal(ctx, conn, dataExec...)
}

func (p PgxAdapter) ExecuteTransactionsWithContext(schema string, ctx context.Context, dataExec ...[]DataExec) error {

	conn, err := p.db.Acquire(ctx)
	if err != nil {
		slog.Error("Fallo conexión db", "error", err)
		return err
	}
	defer conn.Release() // SEGURO: Siempre vuelve al pool al terminar la función

	// Configurar el esquema para ESTA sesión
	if err := p.setSchema(ctx, conn, schema); err != nil {
		return err
	}
	// Iniciar transacción
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// Ejecutar cada grupo de queries usando el motor interno
	for _, group := range dataExec {
		// Pasamos 'tx' como el ejecutor. Si falla, el motor devuelve error.
		if err := p.executeInternal(ctx, tx, group...); err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}
