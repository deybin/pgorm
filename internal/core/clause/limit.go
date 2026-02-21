package clause

import (
	"strconv"
	"strings"
)

type Limit struct {
	Limit  int
	Offset int
}

func (l Limit) Name() string {
	return "LIMIT"
}

func (l Limit) offset() string {
	return "OFFSET"
}

func (l *Limit) Set(limit ...int) {
	switch len(limit) {
	case 1:
		l.Limit = limit[0]
	case 2:
		l.Limit = limit[0]
		l.Offset = limit[1]
	default:
		l.Limit = 0
		l.Offset = 0

	}

}

func (l *Limit) Reset() {
	l.Limit = 0
	l.Offset = 0
}

func (l Limit) Build() string {
	var querySQL strings.Builder

	if l.Limit > 0 {
		querySQL.WriteString(l.Name())
		querySQL.WriteByte(' ')
		querySQL.WriteString(strconv.Itoa(l.Limit))
		querySQL.WriteByte(' ')
		if l.Offset > 0 {
			querySQL.WriteString(l.offset())
			querySQL.WriteByte(' ')
			querySQL.WriteString(strconv.Itoa(l.Offset))
			querySQL.WriteByte(' ')
		}
	}
	return querySQL.String()

}
