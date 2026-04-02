package domain

import (
	"github.com/deybin/pgorm/internal/core/clause"
)

/** guarda la estructura de consulta sql, aparir de aquí se generar la consulta sql */
type Sintaxis struct {
	From_field          clause.From
	Select_field        clause.Select
	Where_field         clause.Where
	Join_field          clause.Join
	Limit_field         clause.Limit
	OrderBy_field       clause.OrderBy
	GroupBy_field       clause.GroupBy
	ArgsLen_field       int
	Args_field          []any
	QueryFull_field     string /** guarda la consulta sql directa en string */
	WorkQueryFull_field bool   /** establece si se va a utilizar una consulta directa mediante queryFull o mediante la estructura true:= se considerara queryFull false:= se considerara  estructura para formar la consulta sql*/
}

func NewSintaxis() *Sintaxis {
	return &Sintaxis{}
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
func (q *Sintaxis) WorkQueryFull(query string, arg ...interface{}) *Sintaxis {
	q.WorkQueryFull_field = true
	q.QueryFull_field = query
	if arg == nil {
		return q
	}
	q.Args_field = append(q.Args_field, arg...)

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
func (q *Sintaxis) From(table string) *Sintaxis {
	q.From_field.Table = table
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
func (q *Sintaxis) Select(campos ...string) *Sintaxis {
	if len(campos) > 0 {
		q.Select_field.Columns = append(q.Select_field.Columns, campos...)
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
func (q *Sintaxis) Where(where string, op clause.OperatorWhere, arg any) *Sintaxis {
	q.Where_field.New(clause.ExpressionFilter{Name: clause.WHERE, Column: where, Operators: op, Args: arg})
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
func (q *Sintaxis) And(and string, op clause.OperatorWhere, arg any) *Sintaxis {
	q.Where_field.Set(clause.ExpressionFilter{Name: clause.AND, Column: and, Operators: op, Args: arg})
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
func (q *Sintaxis) Or(or string, op clause.OperatorWhere, arg any) *Sintaxis {
	q.Where_field.Set(clause.ExpressionFilter{Name: clause.OR, Column: or, Operators: op, Args: arg})
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
func (q *Sintaxis) OrderBy(campos ...string) *Sintaxis {
	q.OrderBy_field.Set(campos...)
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
func (q *Sintaxis) Top(top int) *Sintaxis {
	q.Limit_field.Set(top)
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
func (q *Sintaxis) Limit(limit ...int) *Sintaxis {
	q.Limit_field.Set(limit...)
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
func (q *Sintaxis) GroupBy(group ...string) *Sintaxis {
	q.GroupBy_field.Set(group...)
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
func (q *Sintaxis) Join(tp clause.TypeJoin, table string, on string) *Sintaxis {
	q.Join_field.Expressions = append(q.Join_field.Expressions, clause.ExpressionJoin{
		Type:      tp,
		Table:     table,
		Alias:     "",
		Condition: on,
	})
	return q
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
func (q *Sintaxis) Reset() {
	q.From_field.Reset()
	q.Select_field.Reset()
	q.Where_field.Reset()
	q.Join_field.Reset()
	q.Limit_field.Reset()
	q.OrderBy_field.Reset()
	q.GroupBy_field.Reset()
	q.Args_field = []any{}
	q.QueryFull_field = ""
	q.WorkQueryFull_field = false
}
func (q *Sintaxis) Arguments() []any {
	return q.Args_field
}
