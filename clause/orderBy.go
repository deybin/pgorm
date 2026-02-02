package clause

import (
	"strings"
)

type OrderBy struct {
	col []string
}

func (o OrderBy) Name() string {
	return "ORDER BY"
}

func (o *OrderBy) Set(col ...string) {
	o.col = append(o.col, col...)
}

func (o *OrderBy) Reset() {
	o.col = []string{}
}

func (o OrderBy) Build() string {
	if len(o.col) <= 0 {
		return ""
	}

	var querySQL strings.Builder
	querySQL.WriteString(o.Name())
	querySQL.WriteByte(' ')
	querySQL.WriteString(strings.Join(o.col, ", "))
	querySQL.WriteByte(' ')
	return querySQL.String()
}
