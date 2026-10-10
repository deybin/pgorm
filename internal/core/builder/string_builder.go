package builder

import (
	"errors"
	"strings"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/domain"
)

func BuilderScriptGeneric(ts *domain.Transactions, query string, data ...any) error {
	var sqlExec = make([]adapters.DataExec, 0)
	sqlExec = append(sqlExec, adapters.DataExec{
		Querys: query,
		Values: data,
		Action: ts.Action(),
	})
	// fmt.Println(sqlExec)

	if strings.TrimSpace(query) != "" {
		ts.SetQuery(sqlExec)
		return nil
	} else {
		return errors.New("no existen datos para realizar action")
	}
}
