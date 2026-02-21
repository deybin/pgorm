package builder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/migrator"
)

func BuilderInsertGeneric(ts *domain.Transactions) error {
	table := ts.Schema().Table().Name()
	data := ts.Datos()
	schema := ts.Schema().ParseInsert()
	length := len(data)
	if length > 0 {
		var sqlExec = make([]adapters.DataExec, 0)
		var data_insert []map[string]any

		for _, item := range data {
			preArray, err := migrator.CheckInsertGeneric(schema, item)
			if err == nil {
				data_insert = append(data_insert, preArray)
				var column []string
				var values []string
				var i int
				var valuesExec []any
				char := "$"
				for k, v := range preArray {
					i++
					column = append(column, k)
					values = append(values, fmt.Sprintf("%s%d", char, i))
					valuesExec = append(valuesExec, v)
				}

				sqlPreparate := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, strings.Join(column, ", "), strings.Join(values, ", "))
				sqlExec = append(sqlExec, adapters.DataExec{
					Querys: sqlPreparate,
					Values: valuesExec,
					Action: ts.Action(),
				})
			} else {
				return err
			}
		}
		ts.SetQuery(sqlExec)
		ts.SetData(data_insert)
		return nil
	} else {
		return errors.New("no existen datos para insertar")
	}
}
