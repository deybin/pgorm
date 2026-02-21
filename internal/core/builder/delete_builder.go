package builder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/migrator"
)

func BuilderDeleteGeneric(ts *domain.Transactions, data []migrator.Where) error {
	table := ts.Schema().Table().Name()
	schemas := ts.Schema().ParseUpdate()
	length := len(data)

	if length > 0 {
		var sqlExec = make([]adapters.DataExec, 0)
		// var data_delete []map[string]any

		preArray, err := migrator.CheckWhereGeneric(schemas, data...)
		if err != nil {
			return err
		}

		sqlWherePreparateDelete := ""
		var i uint64

		var valuesExec []interface{}

		char := "$"
		var wheres []string

		for _, v := range preArray {
			i++
			wheres = append(wheres, fmt.Sprintf("%s %s %s %s%d", v.Clause, v.Field, v.Condition, char, i))
			valuesExec = append(valuesExec, v.Value)
		}

		sqlWherePreparateDelete = strings.Join(wheres, " ")

		sqlPreparate := fmt.Sprintf("DELETE FROM %s %s", table, sqlWherePreparateDelete)
		sqlExec = append(sqlExec, adapters.DataExec{Querys: sqlPreparate, Values: valuesExec, Action: ts.Action()})

		ts.SetQuery(sqlExec)
		return nil
	} else {
		return errors.New("no existen datos para eliminar")
	}
}
