package builder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/migrator"
)

func BuilderUpdateGeneric(ts *domain.Transactions) error {
	table := ts.Schema().Table().Name()
	data := ts.Datos()
	schemas := ts.Schema().ParseUpdate()
	length := len(data)

	if length > 0 {
		var sqlExec = make([]adapters.DataExec, 0)
		var data_update []map[string]any
		for _, item := range data {
			v := reflect.ValueOf(item)
			valData := v.FieldByName("Entity")
			valueData := valData.Interface()

			preArray, err := migrator.CheckUpdateGeneric(schemas, valueData)
			if err != nil {
				return err
			}

			if len(preArray) <= 0 {
				continue
			}

			val := v.FieldByName("Conditions")
			value := val.Interface().([]migrator.Where)
			var preArray_where []migrator.Where
			lengthWhere := len(value)
			if lengthWhere > 0 {
				preArray, err := migrator.CheckWhereGeneric(schemas, value...)
				if err != nil {
					return err
				}
				preArray_where = preArray
			}

			data_update = append(data_update, preArray)
			var setters []string

			sqlWherePreparateUpdate := ""
			var i uint64
			var valuesExec []interface{}
			char := "$"
			for k, v := range preArray {
				i++
				setters = append(setters, fmt.Sprintf("%s= %s%d", k, char, i))
				valuesExec = append(valuesExec, v)
			}

			if lengthWhere > 0 {

				var wheres []string
				for _, v := range preArray_where {
					i++
					wheres = append(wheres, fmt.Sprintf("%s %s %s %s%d", v.Clause, v.Field, v.Condition, char, i))
					valuesExec = append(valuesExec, v.Value)
				}

				sqlWherePreparateUpdate = strings.Join(wheres, " ")
				// fmt.Println(sqlWherePreparateUpdate)

			}
			sqlPreparate := fmt.Sprintf("UPDATE %s SET %s %s", table, strings.Join(setters, ", "), sqlWherePreparateUpdate)
			sqlExec = append(sqlExec, adapters.DataExec{
				Querys: sqlPreparate,
				Values: valuesExec,
				Action: ts.Action(),
			})

		}

		if len(data_update) <= 0 {
			return errors.New("al realizar validaciones se filtro datos y se quedo sin información para actualizar")
		}
		ts.SetQuery(sqlExec)
		ts.SetData(data_update)
		ts.SetAction(ts.Action())

		return nil
	} else {
		return errors.New("no existen datos para actualizar")
	}
}
