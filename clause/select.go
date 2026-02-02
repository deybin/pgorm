package clause

import (
	"strings"
)

// Select select
type Select struct {
	Distinct bool
	Columns  []string
	// Expression Expression
}

func (s Select) Name() string {
	return "SELECT "
}

func (s *Select) Reset() {
	s.Distinct = false
	s.Columns = []string{}
}

func (s Select) Build() string {
	var SQL strings.Builder
	SQL.WriteString(s.Name())
	if len(s.Columns) > 0 {
		if s.Distinct {
			SQL.WriteString("DISTINCT ")
		}

		SQL.WriteString(strings.Join(s.Columns, ","))

	} else {
		SQL.WriteByte('*')
	}
	SQL.WriteByte(' ')
	return SQL.String()
}
