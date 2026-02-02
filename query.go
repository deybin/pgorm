package pgorm

import (
	"fmt"
	"log"
	"strings"

	"github.com/deybin/pgorm/clause"
	"github.com/deybin/pgorm/internal"
	"github.com/deybin/pgorm/logger"
	"github.com/deybin/pgorm/schema"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Query struct {
	query         sintaxis
	rowSql        pgx.Rows
	colSql        []string
	conn          *internal.Connection
	err           error
	sessionActiva bool
}

/** guarda la estructura de consulta sql, aparir de aquí se generar la consulta sql */
type sintaxis struct {
	From          clause.From
	Select        clause.Select
	Where         clause.Where
	Join          clause.Join
	Limit         clause.Limit
	OrderBy       clause.OrderBy
	GroupBy       clause.GroupBy
	argsLen       int
	args          []interface{}
	queryFull     string /** guarda la consulta sql directa en string */
	workQueryFull bool   /** establece si se va a utilizar una consulta directa mediante queryFull o mediante la estructura true:= se considerara queryFull false:= se considerara  estructura para formar la consulta sql*/
}

type Config struct {
	Cloud     bool
	Database  string
	Procedure bool
}

func NewQuery(config Config) *Query {
	var err error
	var q Query
	if config.Cloud {
		if q.conn, err = new(internal.Connection).New("").PoolMaster(); err != nil {
			q.err = err
			return &q
		}
	} else {
		if q.conn, err = new(internal.Connection).New(config.Database).NewPool(); err != nil {
			q.err = err
			return &q
		}
	}
	//q.procedure = config.Procedure

	// if _, err = q.conn.GetCnn().Exec(q.conn.Context(), "SET search_path TO public"); err != nil {
	// 	q.err = err
	// 	return q
	// }
	return &q
}

/*
SetQueryString establece una consulta SQL completa y sus argumentos manualmente en el struct Query.

Esta función permite definir una sentencia SQL personalizada que no se construyó mediante los métodos
de composición (como Select, Where, Join, etc.). Es útil cuando se desea ejecutar directamente
una consulta específica con sus parámetros.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetQueryString("SELECT * FROM my_table WHERE email = $1", "ejemplo@correo.com")

Parámetros:
  - query (string): Cadena con la consulta SQL completa.
  - arg (...interface{}): Argumentos de la consulta SQL, pueden ser uno o varios valores.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) WorkQueryFull(query string, arg ...interface{}) *Query {
	q.query.workQueryFull = true
	q.query.queryFull = query
	if arg == nil {
		return q
	}
	q.query.args = append(q.query.args, arg...)

	return q
}

/*
From establece el nombre de la tabla que se utilizará en la consulta SQL.

Esta función define el nombre de la tabla principal sobre la cual se realizarán
las operaciones SQL (SELECT, JOIN, WHERE, etc.). Es esencial establecer esta propiedad
antes de construir la consulta.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.From("my_table").Select("id", "nombre")

Parámetros:
  - table (string): Nombre de la tabla como string.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) From(table schema.Schema) *Query {
	q.query.From.Table = table.Name()
	return q
}

/*
Select establece la cláusula SELECT de la consulta SQL.
Puede especificar una lista de campos como argumentos.
Si no se proporcionan campos, se seleccionan todos (*).

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select().Exec().All()
	result, err := queryBuilder.SetTable("my_table").Select("campo1,campo2").Exec().All()
	result, err := queryBuilder.SetTable("my_table").Select("campo1", "campo2").Exec().All()
	consultaFinal := queryBuilder.String()

Parámetros:
  - campos(...string): Lista de nombres de campos a seleccionar. Si está vacío, se seleccionan todos los campos.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) Select(campos ...string) *Query {
	if len(campos) > 0 {
		q.query.Select.Columns = append(q.query.Select.Columns, campos...)
	}
	return q
}

/*
Where establece la cláusula WHERE de la consulta SQL con una condición y un operador.
La condición puede contener placeholders ($) para argumentos de la consulta.
El operador se utiliza para comparar valores en la condición.
El argumento es el valor que se comparará en la condición.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1, campo2").Where("campo3", "=", valor).Exec().All()

	consultaFinal := queryBuilder.String()

Parámetros:
  - where (string): Campo o expresión de la condición a evaluar.
  - op (OperatorWhere): Operador lógico o de comparación (por ejemplo: "=", "<>", ">", "<", "<=", ">=", "LIKE", "IN", etc.).
  - arg (interface{}): Valor o valores con los que se compara (puede ser simple o slice dependiendo del operador).

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) Where(where string, op clause.OperatorWhere, arg interface{}) *Query {
	q.query.Where.Set(clause.ExpressionFilter{Name: q.query.Where.Name(), Column: where, Operators: op, Args: arg})
	return q
}

/*
And añade una cláusula AND adicional a la cláusula WHERE existente de la consulta SQL.

Esta función permite agregar condiciones adicionales a la cláusula WHERE ya existente.
Si no se ha definido una cláusula WHERE previa, no se realiza ninguna acción.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1, campo2").Where("campo3", "=", valor).And("campo4", ">", otroValor).Exec().All()

	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - and (string): Campo o expresión adicional para la condición.
  - op (OperatorWhere): Operador lógico de comparación (por ejemplo: "=", "<>", ">", "<", "<=", ">=", "LIKE", "IN", etc.).
  - arg (interface{}): Valor a comparar (puede ser simple o slice dependiendo del operador).

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) And(and string, op clause.OperatorWhere, arg interface{}) *Query {
	q.query.Where.Set(clause.ExpressionFilter{Name: q.query.Where.And(), Column: and, Operators: op, Args: arg})
	return q
}

/*
Or añade una cláusula OR adicional a la cláusula WHERE existente de la consulta SQL.

Esta función permite agregar condiciones adicionales a la cláusula WHERE ya existente.
Si no se ha definido una cláusula WHERE previa, no se realiza ninguna acción.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1, campo2").Where("campo3", "=", valor).Or("campo4", ">", otroValor).Exec().All()

	consultaFinal := queryBuilder.GetQuery()

Parámetros:
  - or (string): Campo o expresión adicional para la condición.
  - op (OperatorWhere): Operador lógico de comparación (por ejemplo: "=", "<>", ">", "<", "<=", ">=", "LIKE", "IN", etc.).
  - arg (interface{}): Valor a comparar (puede ser simple o slice dependiendo del operador).

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) Or(or string, op clause.OperatorWhere, arg interface{}) *Query {
	q.query.Where.Set(clause.ExpressionFilter{Name: q.query.Where.Or(), Column: or, Operators: op, Args: arg})
	return q
}

/*
OrderBy establece la cláusula ORDER BY de la consulta SQL.
Permite ordenar los resultados de la consulta según uno o más campos especificados.

Cada campo puede ir acompañado de la dirección de ordenamiento (`ASC` o `DESC`).

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1", "campo2").Where("campo3", "=", valor).OrderBy("campo4 DESC", "campo5 ASC").Exec().All()
	consultaFinal := queryBuilder.String()

Parámetros:
  - campos(...string): Lista de nombres de columnas con o sin dirección de ordenamiento.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) OrderBy(campos ...string) *Query {
	q.query.OrderBy.Set(campos...)
	return q
}

/*
Top establece la cláusula LIMIT de la consulta SQL para limitar el número de filas devueltas.

Esta función es útil cuando solo deseas recuperar una cantidad específica de registros,
por ejemplo, para paginación o para obtener los primeros resultados de una tabla ordenada.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1", "campo2").Top(10).Exec().All()

	consultaFinal := queryBuilder.String()

Parámetros:
  - top (int): Número máximo de filas que se desean recuperar.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) Top(top int) *Query {
	q.query.Limit.Set(top)
	return q
}

/*
Limit establece la cláusula LIMIT (y opcionalmente OFFSET) de la consulta SQL.

Esta función permite limitar la cantidad de resultados devueltos por la consulta
y, si se desea, omitir un número determinado de registros iniciales (offset), útil para paginación.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})

	// Limita la consulta a 10 filas
	result, err := queryBuilder.SetTable("my_table").Select().Where("campo3", "=", valor).Limit(10).Exec().All()

	// Limita la consulta a 10 filas, omitiendo las primeras 3
	result, err := queryBuilder.SetTable("my_table").Select("campo1, campo2").Where("campo3", "=", valor).Limit(10, 3).Exec().All()

	consultaFinal := queryBuilder.String()

Parámetros:
  - limit (...int): Lista de uno o dos enteros.
  - Si se proporciona un solo entero, aplica solo el límite de filas.
  - Si se proporcionan dos, el primero es el límite y el segundo el offset.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) Limit(limit ...int) *Query {
	q.query.Limit.Set(limit...)
	return q
}

/*
GroupBy establece la cláusula GROUP BY de la consulta SQL.
Se utiliza para agrupar los resultados de una consulta por uno o más campos especificados.

Esto es útil especialmente cuando se aplican funciones de agregación como COUNT, SUM, AVG, etc.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1", "COUNT(campo2)").GroupBy("campo1").Exec().All()
	consultaFinal := queryBuilder.String()

Parámetros:
  - group (...string): Lista de nombres de campos a utilizar en la cláusula GROUP BY.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) GroupBy(group ...string) *Query {
	q.query.GroupBy.Set(group...)
	return q
}

/*
Join añade una cláusula JOIN a la consulta SQL.

Permite establecer una relación entre la tabla principal y otra tabla especificando
el tipo de unión (INNER, LEFT, RIGHT, FULL), la tabla a unir y la condición de emparejamiento (ON).

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("tabla_principal.columna1, tabla_secundaria.columna2").
		Join(pgorm.INNER, "tabla_secundaria", "tabla_principal.id = tabla_secundaria.id").
		Where("tabla_principal.columna3", "=", valor).Exec().All()

	consultaFinal := queryBuilder.String()

Parámetros:
  - tp (TypeJoin): Tipo de unión (TypeJoin). Puede ser pgorm.INNER, pgorm.LEFT, pgorm.RIGHT o pgorm.FULL.
  - table (string): Nombre de la tabla a unir.
  - on (string): Condición ON que define cómo se relacionan las tablas.

Devuelve:
  - Un puntero al struct Query actualizado para permitir el encadenamiento de métodos.
*/
func (q *Query) Join(tp clause.TypeJoin, table schema.Schema, on string) *Query {
	q.query.Join.Expressions = append(q.query.Join.Expressions, clause.ExpressionJoin{
		Type:      tp,
		Table:     table.Name(),
		Alias:     "",
		Condition: on,
	})
	return q
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
func (q *Query) Exec() *Query {
	if q.err != nil {
		return q
	}

	q.sessionActiva = false
	queryString := q.build()
	// fmt.Println("query:", queryString)

	rows, err := q.conn.Pool().Query(q.conn.Context(), queryString, q.query.args...)
	if err != nil {
		q.err = err
		return q
	}

	q.rowSql = rows
	q.keyFieldName()
	return q

}

/*
ExecCtx ejecuta la consulta SQL construida utilizando `pgx` y almacena los resultados en la estructura `Query`.

A diferencia de otras variantes como `Exec`, esta función **no cierra la conexión con la base de datos**, permitiendo que
se realicen múltiples consultas consecutivas sobre la misma conexión. Por esta razón, es responsabilidad del desarrollador
cerrar la conexión manualmente mediante `pgorm.Query.Close()` cuando ya no se requiera.

Esta función ejecuta la consulta previamente construida mediante los métodos del builder (`Select`, `Where`, `SetTable`, etc.)
o mediante una consulta SQL personalizada definida con `SetQueryString`. Soporta tanto consultas estándar (`SELECT`)
como procedimientos (`CALL`, `PROCEDURE`).

- Si se trata de una consulta `SELECT`, almacena los resultados (`pgx.Rows`) en `q.rowSql` y los nombres de las columnas en `q.colSql`.
- Si la consulta no retorna resultados (`INSERT`, `UPDATE`, `DELETE`, etc.), la ejecuta directamente y captura cualquier error.

Utiliza el contexto (`q.conn.Context()`) asociado a la conexión para soportar cancelación y timeout.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetTable("my_table").
	      Select("id, name").
	      Where("status", "=", "active").
	      ExecCtx()

Nota:
  - La conexión **no se cierra automáticamente**. Debes cerrarla manualmente con `queryBuilder.Close()` cuando finalices el uso del objeto `Query`.

Retorna:
  - Un puntero a `Query`, con los resultados cargados o el error ocurrido.
*/
func (q *Query) ExecCtx() *Query {
	if q.err != nil {
		return q
	}
	q.sessionActiva = true

	queryString := q.build()
	// fmt.Println("query:", queryString)

	rows, err := q.conn.Pool().Query(q.conn.Context(), queryString, q.query.args...)
	if err != nil {
		q.err = err
		return q
	}

	q.rowSql = rows
	q.keyFieldName()
	return q

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
func (q *Query) Procedure() error {
	if q.err != nil {
		return q.err
	}

	q.sessionActiva = false

	defer q.conn.Close()
	queryString := q.build()
	//fmt.Println("query:", queryString)
	fmt.Println(len(q.query.args), q.query.args)

	if _, err := q.conn.Pool().Exec(q.conn.Context(), queryString, q.query.args...); err != nil {
		q.err = err
		return err
	}
	return nil

}

/*
ProcedureCtx ejecuta una consulta de tipo procedimiento (por ejemplo, `CALL` o funciones sin retorno de datos) utilizando `pgx`,
manteniendo activa la sesión para continuar ejecutando más consultas sobre la misma conexión.

Esta función es útil cuando necesitas ejecutar procedimientos almacenados u operaciones SQL sin retorno de resultados,
pero planeas seguir utilizando la conexión activa después de su ejecución.

A diferencia de `Procedure`, esta función **no cierra automáticamente la conexión** (`q.conn`), por lo que es responsabilidad
del desarrollador cerrarla manualmente con `pgorm.Query.Close()` al finalizar el uso.

Comportamiento:
- Si existe un error previo en el estado del `Query` (`q.err`), la función lo retorna inmediatamente.
- Marca la sesión como activa (`q.sessionActiva = true`) para indicar que la conexión sigue disponible.
- Ejecuta la consulta generada por `getQuery()` con los argumentos acumulados (`q.args`).
- Si ocurre un error durante la ejecución, lo guarda en `q.err` y lo retorna.

Ejemplo de uso:

	query := new(pgorm.Query).New(pgorm.QConfig{Database: "my_db"})
	query.SetQueryString("CALL procesar_datos($1)").SetArgs(123).ProcedureCtx()
	// ... puedes seguir usando query.conn mientras la sesión esté activa

Nota:
  - La conexión **debe cerrarse manualmente** con `query.Close()` cuando ya no se necesite.

Retorna:
  - `nil` si la ejecución fue exitosa,
  - Un `error` si falló la ejecución o ya existía un error previo en la estructura.
*/
func (q *Query) ProcedureCtx() error {
	if q.err != nil {
		return q.err
	}

	q.sessionActiva = true
	queryString := q.build()
	// fmt.Println("query:", queryString)

	if _, err := q.conn.Pool().Exec(q.conn.Context(), queryString, q.query.args...); err != nil {
		q.err = err
		return err
	}
	return nil

}

/*
One recupera un solo resultado de la consulta SQL y lo devuelve como un mapa[string]interface{}.

Esta función se utiliza para leer una única fila del conjunto de resultados almacenado en `q.rowSql`,
que debe haber sido previamente cargado mediante el método `Exec()`.

Devuelve un mapa donde:
- Las claves son los nombres de las columnas devueltas por la consulta.
- Los valores son los correspondientes valores de cada columna en la primera fila.

Si no hay filas o ocurre un error, se retorna un error o un mapa vacío.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1, campo2").Where("campo3", "=", valor).Exec().One()

Devuelve:
  - Un `map[string]interface{}` con los datos de la primera fila del resultado.
  - Un error, si no se encuentra una fila o ocurre algún problema al escanear los datos.
*/
func (q *Query) One() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if q.err != nil {
		log.Printf("check failed: %v", q.err)
		return m, logger.ManagerErrors{}.SqlQuery(q.err)
	}

	if !q.sessionActiva {
		defer q.conn.Close()
	}

	defer q.rowSql.Close()
	fieldDescs := q.rowSql.FieldDescriptions()
	if q.rowSql.Next() {
		row, err := q.builderResult()
		if err != nil {
			log.Printf("check failed: %v", err)
			return map[string]interface{}{}, logger.ManagerErrors{}.SqlQuery(err)
		}
		row = q.normalizeRow(row, fieldDescs)
		m = row

	}

	return m, nil
}

/*
Text recupera el valor de una columna específica de la primera fila de resultados de la consulta SQL.

Esta función es útil cuando se necesita obtener directamente un valor puntual de la primera fila,
sin necesidad de recorrer todo el conjunto de resultados.

Debe ser llamada después de ejecutar la consulta con `Exec()`.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetTable("my_table").Select("campo1, campo2").Where("campo3", "=", valor).Exec()
	valor, err := queryBuilder.Text("nombreColumna")

Parámetros:
  - columna (string): Nombre de la columna cuyo valor se desea recuperar.

Devuelve:
  - El valor de la columna especificada (`interface{}`).
  - Un error si no se encontró la columna, si no hay filas, o si ocurrió algún problema al leer el valor.
*/
func (q *Query) Text(columna string) (interface{}, error) {

	if q.err != nil {
		log.Printf("check failed: %v", q.err)
		return nil, logger.ManagerErrors{}.SqlQuery(q.err)
	}

	if !q.sessionActiva {
		defer q.conn.Close()
	}

	defer q.rowSql.Close()
	m := make(map[string]interface{})
	fieldDescs := q.rowSql.FieldDescriptions()
	if q.rowSql.Next() {
		row, err := q.builderResult()
		if err != nil {
			log.Printf("check failed: %v", err)
			return nil, logger.ManagerErrors{}.SqlQuery(err)
		}
		row = q.normalizeRow(row, fieldDescs)
		m = row

	}

	return m[columna], nil
}

/*
All recupera todas las filas de resultados de la consulta SQL y las devuelve como una lista de mapas.

Esta función debe ser utilizada después de ejecutar una consulta con `Exec()`.
Cada resultado se representa como un mapa, donde las claves son los nombres de las columnas
y los valores son los valores correspondientes de esa fila.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	result, err := queryBuilder.SetTable("my_table").Select("campo1", "campo2").Where("campo3", "=", valor).Exec().All()

Devuelve:
  - Una lista de mapas (`[]map[string]interface{}`), donde cada mapa representa una fila de resultados.
  - Un error, si ocurre alguno durante el procesamiento de las filas.
*/
func (q *Query) All() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	if q.err != nil {
		log.Printf("check failed: %v", q.err)
		return result, logger.ManagerErrors{}.SqlQuery(q.err)
	}

	if !q.sessionActiva {
		defer q.conn.Close()
	}

	defer q.rowSql.Close()

	fieldDescs := q.rowSql.FieldDescriptions()
	for q.rowSql.Next() {
		row, err := q.builderResult()
		if err != nil {
			log.Printf("check failed: %v", err)
			return []map[string]interface{}{}, logger.ManagerErrors{}.SqlQuery(err)
		}
		row = q.normalizeRow(row, fieldDescs)
		result = append(result, row)
	}

	return result, nil
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
func (q *Query) builderResult() (map[string]interface{}, error) {

	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.

	columns := make([]interface{}, len(q.colSql))
	columnPointers := make([]interface{}, len(q.colSql))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := q.rowSql.Scan(columnPointers...); err != nil {
		log.Printf("check failed: %v", err)
		return map[string]interface{}{}, logger.ManagerErrors{}.SqlQuery(err)
	}

	//Crea nuestro mapa y recupera el valor de cada columna del segmento de punteros, almacenándolo en el mapa con el nombre de la columna como clave.
	m := make(map[string]interface{})
	for i, colName := range q.colSql {
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
func (q *Query) keyFieldName() {
	fieldDescription := q.rowSql.FieldDescriptions()
	q.colSql = make([]string, len(fieldDescription))
	for i, fd := range fieldDescription {
		q.colSql[i] = string(fd.Name)
	}
}

/*
build devuelve la cadena completa de la consulta SQL construida con los métodos del struct Querys.

Esta función compone y retorna la consulta SQL generada hasta el momento, integrando las cláusulas
SELECT, JOIN, WHERE, GROUP BY, ORDER BY, LIMIT o TOP, según hayan sido definidas previamente.

Es útil para depurar o inspeccionar la consulta antes de ejecutarla.

Devuelve:
  - Una cadena (`string`) que representa la consulta SQL completa construida.
*/
func (q *Query) build() string {
	var querySql strings.Builder
	// var queryString string
	if !q.query.workQueryFull {
		querySql.WriteString(q.query.Select.Build())
		querySql.WriteString(q.query.From.Build())
		querySql.WriteString(q.query.Join.Build())
		/** aplicando Where : where ,and ,or ,in, between ,not in ,not between*/
		querySql.WriteString(q.query.Where.Build())
		q.query.args = q.query.Where.FindArguments()
		q.query.argsLen = q.query.Where.FindArgumentsLen()

		/** aplicando Group by*/
		querySql.WriteString(q.query.GroupBy.Build())

		/** aplicando order by  */
		querySql.WriteString(q.query.OrderBy.Build())
		/** aplicando Top y LImit  */
		querySql.WriteString(q.query.Limit.Build())
	} else {
		querySql.WriteString(q.query.queryFull)
	}

	return querySql.String()
}

/*
String devuelve la consulta SQL construida como una cadena.

Esta función implementa la interfaz `Stringer` para el struct Query,
permitiendo que la consulta SQL construida se represente directamente como una cadena,
por ejemplo al utilizar fmt.Println(q) o log.Print(q).

Devuelve:
  - Una cadena (`string`) que representa la consulta SQL construida mediante los métodos del struct Query.
*/
func (q *Query) String() string {
	return q.build()
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
	return q.err
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
	q.rowSql.Close()
	q.conn.Close()
}

/*
Reset reinicia la configuración de la consulta SQL en el struct Query.

Esta función limpia todos los componentes relacionados con la construcción de la consulta SQL,
como las cláusulas SELECT, WHERE, JOIN, etc., así como los argumentos asociados.
Es útil cuando se desea reutilizar el mismo struct `Query` para construir una nueva consulta desde cero.

Ejemplo de uso:

	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "my_database"})
	queryBuilder.SetTable("my_table").Select("campo1").Where("campo2", "=", valor).Exec()
	queryBuilder.Reset() // Limpia la consulta actual
	queryBuilder.SetTable("my_table").Select("otro_campo").Where("otro_campo", "=", nuevoValor).Exec()

No devuelve ningún valor.
*/
func (q *Query) Reset() {
	q.query = sintaxis{}
}

// normalizeRow convierte UUIDs (OID 2950) a string
func (q *Query) normalizeRow(row map[string]interface{}, fieldDescs []pgconn.FieldDescription) map[string]interface{} {
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
