package clause

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/** operaciones utilizadas con la sentencia WHERE*/
type OperatorWhere string

const (
	I           OperatorWhere = "="
	D           OperatorWhere = "<>"
	MY          OperatorWhere = ">"
	MYI         OperatorWhere = ">="
	MN          OperatorWhere = "<"
	MNI         OperatorWhere = "<="
	LIKE        OperatorWhere = "LIKE"
	IN          OperatorWhere = "IN"
	NOT_IN      OperatorWhere = "NOT IN"
	BETWEEN     OperatorWhere = "BETWEEN"
	NOT_BETWEEN OperatorWhere = "NOT BETWEEN"
)

// Where where clause
type Where struct {
	Expressions  []ExpressionFilter
	Arguments    []interface{}
	ArgumentsLen int
}

type ExpressionFilter struct {
	Name      string
	Column    string
	Operators OperatorWhere
	Args      interface{}
}

// Name where clause name
func (w Where) Name() string {
	return "WHERE"
}

// Name where clause name
func (w Where) And() string {
	return "AND"
}

// Name where clause name
func (w Where) Or() string {
	return "OR"
}

// Name where clause name
func (w Where) Not() string {
	return "NOT"
}

func (w *Where) Set(expression ExpressionFilter) {
	w.Expressions = append(w.Expressions, expression)
}

func (w Where) FindArguments() []interface{} {
	return w.Arguments
}
func (w Where) FindArgumentsLen() int {
	return w.ArgumentsLen
}

func (w *Where) Reset() {
	w.Expressions = []ExpressionFilter{}
	w.Arguments = []interface{}{}
	w.ArgumentsLen = 0
}

/*
getSintaxisFilter genera la sintaxis SQL adecuada para filtros aplicados en cláusulas WHERE
y gestiona los argumentos necesarios para operadores como IN, NOT IN, BETWEEN y NOT BETWEEN.

Esta función se encarga de construir la representación textual correcta del filtro en función del operador proporcionado,
y añade los argumentos correspondientes a la lista de parámetros de la consulta.

Parámetros:
  - op: Operador de filtro (OperatorWhere) que se aplicará (por ejemplo, "=", "IN", "BETWEEN", etc.).
  - arg: Valor o conjunto de valores utilizados en el filtro. Puede ser un valor único o un slice (por ejemplo, []interface{}).

Devuelve:
  - Una cadena (`string`) que representa la sintaxis adecuada para el filtro SQL.
  - Un error (`error`) si ocurre alguno durante la generación de la sintaxis o el manejo de los argumentos.
*/
func (w *Where) buildExp(expr ExpressionFilter) (string, error) {
	// var argString string
	var SQL strings.Builder

	SQL.WriteString(expr.Name)
	SQL.WriteByte(' ')

	switch expr.Operators {
	case IN, NOT_IN:
		if reflect.TypeOf(expr.Args).String() != "[]interface {}" {
			return "", errors.New("tipo de dato incorrecto para filtrado IN")
		}
		if len(expr.Args.([]interface{})) <= 0 {
			return "", errors.New("valor vació para filtrado IN")
		}

		SQL.WriteString(expr.Column)
		SQL.WriteByte(' ')
		SQL.WriteString(string(expr.Operators))
		SQL.WriteByte(' ')

		arrayArgsSql := make([]string, 0)
		for _, v := range expr.Args.([]interface{}) {
			arrayArgsSql = append(arrayArgsSql, fmt.Sprintf("$%d", w.ArgumentsLen))
			w.Arguments = append(w.Arguments, v)
			w.ArgumentsLen++
		}
		SQL.WriteByte('(')
		SQL.WriteString(strings.Join(arrayArgsSql, ", "))
		SQL.WriteByte(')')
	case BETWEEN, NOT_BETWEEN:
		if reflect.TypeOf(expr.Args).String() != "[]interface {}" {
			return "", errors.New("tipo de dato incorrecto para filtrado BETWEEN")
		}
		if len(expr.Args.([]interface{})) < 2 {
			return "", errors.New("valor vació o bien valores incompletos para filtrado BETWEEN")
		}

		SQL.WriteString(expr.Column)
		SQL.WriteByte(' ')
		SQL.WriteString(string(expr.Operators))
		SQL.WriteByte(' ')
		SQL.WriteByte('$')
		SQL.WriteString(strconv.Itoa(w.ArgumentsLen))
		SQL.WriteString(" AND ")
		// argString = fmt.Sprintf("$%d AND ", q.argsLen)
		w.Arguments = append(w.Arguments, expr.Args.([]interface{})[0])
		w.ArgumentsLen++
		SQL.WriteByte('$')
		SQL.WriteString(strconv.Itoa(w.ArgumentsLen))
		// argString += fmt.Sprintf("$%d", q.argsLen)
		w.Arguments = append(w.Arguments, expr.Args.([]interface{})[1])
		w.ArgumentsLen++
	default:
		SQL.WriteString(expr.Column)
		SQL.WriteByte(' ')
		SQL.WriteString(string(expr.Operators))
		SQL.WriteByte(' ')
		SQL.WriteByte('$')
		SQL.WriteString(strconv.Itoa(w.ArgumentsLen))
		w.Arguments = append(w.Arguments, expr.Args)
		w.ArgumentsLen++
	}

	// fmt.Println(SQL.String())
	return SQL.String(), nil
}

func (w *Where) Build() string {
	var SQL strings.Builder
	w.ArgumentsLen = 1
	for _, v := range w.Expressions {
		script, _ := w.buildExp(v)
		// fmt.Println("name:=")
		SQL.WriteString(script)
		SQL.WriteByte(' ')
	}
	// fmt.Println(w.Arguments)

	return SQL.String()
}
