package clause

import (
	"strings"
)

type GroupBy struct {
	col []string
}

func (g GroupBy) Name() string {
	return "GROUP BY"
}

func (g *GroupBy) Set(col ...string) {
	g.col = append(g.col, col...)
}

func (g *GroupBy) Reset() {
	g.col = []string{}
}
func (g GroupBy) Build() string {
	if len(g.col) <= 0 {
		return ""
	}

	var querySQL strings.Builder
	querySQL.WriteString(g.Name())
	querySQL.WriteByte(' ')
	querySQL.WriteString(strings.Join(g.col, ", "))
	querySQL.WriteByte(' ')
	return querySQL.String()
}
