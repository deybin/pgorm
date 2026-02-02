package clause

import (
	"strings"
)

type TypeJoin string

const (
	INNER TypeJoin = "INNER JOIN"
	LEFT  TypeJoin = "LEFT JOIN"
	RIGHT TypeJoin = "RIGHT JOIN"
	FULL  TypeJoin = "FULL OUTER JOIN"
)

type Join struct {
	Expressions []ExpressionJoin
}

type ExpressionJoin struct {
	Type      TypeJoin
	Table     string
	Alias     string
	Condition string
}

func (j Join) Name() string {
	return "JOIN"
}
func (j *Join) Set(expr ExpressionJoin) {
	j.Expressions = append(j.Expressions, expr)
}

func (j Join) Reset() {
	j.Expressions = []ExpressionJoin{}
}

func (j Join) Build() string {
	var querySQL strings.Builder
	for _, v := range j.Expressions {
		querySQL.WriteString(string(v.Type))
		querySQL.WriteByte(' ')
		querySQL.WriteString(v.Table)
		querySQL.WriteString(" ON ")
		querySQL.WriteString(v.Condition)
		querySQL.WriteByte(' ')
	}
	return querySQL.String()
}
