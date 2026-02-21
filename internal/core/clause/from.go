package clause

import (
	"strings"
)

type From struct {
	Table string
}

func (f From) Name() string {
	return "FROM "
}

func (f *From) Reset() {
	f.Table = ""
}

func (f From) Build() string {
	var SQL strings.Builder
	SQL.WriteString(f.Name())
	name := f.Table
	SQL.WriteString(name)
	SQL.WriteByte(' ')
	return SQL.String()
}
