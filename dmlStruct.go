package pgorm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/deybin/pgorm/internal"
	"github.com/deybin/pgorm/logger"
	"github.com/deybin/pgorm/schema"
	"github.com/jackc/pgx/v5"
)

type Transactions struct {
	ob     []schema.Models          //datos para observación
	data   []map[string]interface{} //datos para insertar o actualizar o eliminar
	query  []map[string]interface{}
	schema schema.Schema
	action string
	errors []string
}

type SqlExecSingles struct {
	transaction Transactions
	errors      []string
}

type SqlExecMultiples struct {
	tx          pgx.Tx
	database    string
	transaction []*Transactions
}

/*
New crea una nueva instancia de SqlExecSingle con el esquema y los datos proporcionados.

	Parámetros
		* s {Schema}: esquema de la tabla
		* datos {[]map[string]interface{}}: datos a insertar, actualizar o eliminar

	Return
		- (*SqlExecSingle) retorna  puntero *SqlExecSingle struct
*/
func NewSqlExecSingles(s schema.Schema, datos ...schema.Models) *SqlExecSingles {
	return &SqlExecSingles{transaction: Transactions{ob: datos, schema: s}}
}

func (sq *SqlExecSingles) Transaction() *Transactions {
	return &sq.transaction
}

func (sq *SqlExecSingles) Schema() schema.Schema {
	return sq.transaction.schema
}

func (sq *SqlExecSingles) Datos() []schema.Models {
	return sq.transaction.ob
}

/*
Ejecuta el query

	Return
		- returns {error}: retorna errores ocurridos durante la ejecución
*/
func (sq *SqlExecSingles) Exec(database string, params ...bool) error {
	cnn, err := new(internal.Connection).New(database).NewPool()
	if err != nil {
		return err
	}

	cross := false
	if len(params) == 1 {
		cross = params[0]
	}
	dataExec := sq.transaction.query
	defer cnn.Close()
	for _, item := range dataExec {
		sqlPre := item["sqlPreparate"].(string)
		// fmt.Println("PREPARED: ", sqlPre)
		if cross {
			if sq.transaction.action == "UPDATE" {
				sqlPre = Query_Cross_Update(sqlPre)
			}
		}

		// fmt.Println("PREPARED: ", sqlPre)
		valuesExec := item["valuesExec"].([]interface{})

		if _, err_exec := cnn.Pool().Exec(cnn.Context(), sqlPre, valuesExec...); err_exec != nil {
			return fmt.Errorf("error sql %s: %s", sq.transaction.action, logger.ManagerErrors{}.SqlCrud(err_exec, sq.Schema().Name()).Error())
		}
	}
	return nil
}

/*******************************Crud Transactions************************************/
/*
Valida los datos para insertar y crea el query para insertar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *Transactions) Insert() error {
	sqlExec, data_insert, err := sq._insert(sq.schema.Name(), sq.ob, sq.schema.GetSchemaInsert())
	if err != nil {
		return errors.New(strings.Join(sq.errors, "; "))
	}
	sq.query = sqlExec
	sq.data = data_insert
	sq.action = "INSERT"
	return nil
}

/*
Valida los datos para actualizar y crea el query para actualizar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *Transactions) Update() error {
	sqlExec, data_update, err := sq._update(sq.schema.Name(), sq.ob, sq.schema.GetSchemaUpdate())
	if err != nil {
		return errors.New(strings.Join(sq.errors, "; "))
	}
	sq.query = sqlExec
	sq.data = data_update
	sq.action = "UPDATE"
	return nil
}

/*
Valida los datos para Eliminar y crea el query para Eliminar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *Transactions) Delete(dataDelete ...schema.Where) error {
	sqlExec, err := sq._delete(sq.schema.Name(), dataDelete, sq.schema.GetSchemaDelete())
	if err != nil {
		return errors.New(strings.Join(sq.errors, "; "))
	}
	sq.query = sqlExec
	sq.action = "DELETE"
	return nil
}

func (sq *Transactions) _insert(table string, data []schema.Models, schema []schema.Fields) ([]map[string]any, []map[string]any, error) {
	length := len(data)
	if length > 0 {
		var sqlExec = make([]map[string]any, 0)
		var data_insert []map[string]any

		for _, item := range data {
			preArray, err := sq._checkInsertSchema(schema, item)
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
				sqlExec = append(sqlExec, map[string]any{
					"sqlPreparate": sqlPreparate,
					"valuesExec":   valuesExec,
				})
			} else {
				return nil, nil, err
			}
		}
		return sqlExec, data_insert, nil
	} else {
		return nil, nil, errors.New("no existen datos para insertar")
	}
}

func (sq *Transactions) _update(table string, data []schema.Models, schemas []schema.Fields) ([]map[string]any, []map[string]any, error) {
	length := len(data)

	if length > 0 {
		var sqlExec = make([]map[string]any, 0)
		var data_update []map[string]any
		for _, item := range data {
			v := reflect.ValueOf(item)
			valData := v.FieldByName("Data")
			valueData := valData.Interface()

			preArray, err := sq._checkUpdateSchema(schemas, valueData)
			if err != nil {
				return nil, nil, err
			}

			if len(preArray) <= 0 {
				continue
			}

			val := v.FieldByName("Conditions")
			value := val.Interface().([]schema.Where)
			var preArray_where []schema.Where
			lengthWhere := len(value)
			if lengthWhere > 0 {
				preArray, err := sq._checkWhere(schemas, value...)
				if err != nil {
					return nil, nil, err
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
			sqlExec = append(sqlExec, map[string]interface{}{
				"sqlPreparate": sqlPreparate,
				"valuesExec":   valuesExec,
			})

		}

		if len(data_update) <= 0 {
			return nil, nil, errors.New("al realizar validaciones se filtro datos y se quedo sin información para actualizar")
		}
		return sqlExec, data_update, nil
	} else {
		return nil, nil, errors.New("no existen datos para actualizar")
	}
}

func (sq *Transactions) _delete(table string, data []schema.Where, schemas []schema.Fields) ([]map[string]any, error) {
	length := len(data)

	if length > 0 {
		var sqlExec = make([]map[string]any, 0)
		// var data_delete []map[string]any

		preArray, err := sq._checkWhere(schemas, data...)
		if err != nil {
			return nil, err
		}

		var lineSqlExec = make(map[string]any, 2)
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
		lineSqlExec["sqlPreparate"] = sqlPreparate
		lineSqlExec["valuesExec"] = valuesExec
		sqlExec = append(sqlExec, lineSqlExec)

		return sqlExec, nil
	} else {
		return nil, errors.New("no existen datos para eliminar")
	}
}

func (sq *Transactions) _checkInsertSchema(schema []schema.Fields, tabla_map schema.Models) (map[string]any, error) {

	data := make(map[string]any)
	for _, item := range schema {
		v := reflect.ValueOf(tabla_map)
		// fmt.Println("ERRName:", item.NameOriginal)
		val := v.FieldByName(item.NameOriginal)
		// fmt.Println("ERR:", val)
		isNil := val.IsZero()
		// fmt.Println(isNil, item.NameOriginal)
		defaultIsNil := item.Default == nil
		if !isNil {
			if _, ok := val.Interface().(string); ok {
				if strings.TrimSpace(val.String()) == "" {
					isNil = true
				}
			}
		}

		if !isNil {
			val, err := sq.validaciones(item, val.Interface())
			if err == nil {
				data[item.Name] = val
			} else {
				sq.errors = append(sq.errors, fmt.Sprintf("Se encontró fallas al validar el campo %s \n %s\n", item.Description, err.Error()))
			}
		} else {
			if !defaultIsNil {
				data[item.Name] = item.Default
			} else {
				if item.Required {
					sq.errors = append(sq.errors, fmt.Sprintf("El campo %s es Requerido\n", item.Description))
				}
			}
		}

	}
	if len(sq.errors) > 0 {
		return nil, errors.New("No paso validaciones")
	} else {
		return data, nil
	}
}

func (sq *Transactions) _checkUpdateSchema(schemas []schema.Fields, tabla_map any) (map[string]any, error) {
	data := make(map[string]any)
	for _, item := range schemas {
		if item.NameOriginal == "Conditions" {
			continue
		}
		v := reflect.ValueOf(tabla_map)
		val := v.FieldByName(item.NameOriginal)
		isNil := val.IsZero()
		// fmt.Println(isNil, ":", val.Interface(), ":", item.NameOriginal)
		if !isNil {
			if item.Update {
				value := val.Interface()

				if x, ok := value.(string); ok {
					if strings.TrimSpace(x) == "" {
						if !item.Empty {
							sq.errors = append(sq.errors, fmt.Sprintf("El campo %s no puede estar vació\n", item.Description))
						}
					}
				}

				var val interface{}
				val, err := sq.validaciones(item, value)

				if err == nil {
					keyName := item.Name

					switch item.ArithmeticOperations {
					case schema.Sum:
						keyName = fmt.Sprintf("ADD_%s_SUMA", item.Name)
					case schema.Subtraction:
						keyName = fmt.Sprintf("ADD_%s_SUBTRACTION", item.Name)
					case schema.Multiply:
						keyName = fmt.Sprintf("ADD_%s_MULTIPLY", item.Name)
					case schema.Divide:
						keyName = fmt.Sprintf("ADD_%s_DIVIDE", item.Name)
					}

					data[keyName] = val
				} else {
					sq.errors = append(sq.errors, fmt.Sprintf("Se encontró fallas al validar el campo %s \n %s\n", item.Description, err.Error()))
				}
			} else {
				sq.errors = append(sq.errors, fmt.Sprintf("El campo %s no puede ser modificado\n", item.Description))
			}
		}
	}

	if len(sq.errors) > 0 {
		return nil, errors.New("No paso validaciones")
	} else {
		return data, nil
	}
}

func (sq *Transactions) _checkWhere(schemas []schema.Fields, table_where ...schema.Where) ([]schema.Where, error) {

	usePrimaryKey := false
	fieldNotExistLent := 0
	for _, item := range table_where {
		fieldExist := false
		for _, v := range schemas {
			// fmt.Println(item.Field, " : ", v.Name)
			if item.Field == v.Name {
				if v.Where || v.PrimaryKey {
					if v.PrimaryKey {
						usePrimaryKey = true
					}
					if x, ok := item.Value.(string); ok {
						if strings.TrimSpace(x) == "" {
							sq.errors = append(sq.errors, fmt.Sprintf("El campo %s no puede ser utilizado de esta forma\n", v.Description))
						}

					}

				} else {
					sq.errors = append(sq.errors, fmt.Sprintf("El campo %s no puede ser utilizado de esta forma\n", v.Description))
				}
				fieldExist = true
				continue
			}

		}
		if !fieldExist {
			fieldNotExistLent++
		}

	}

	if fieldNotExistLent > 0 {
		sq.errors = append(sq.errors, "Uno o más campos enviados no son válidos.")
	}

	existePrimaryKey := slices.ContainsFunc(schemas, func(u schema.Fields) bool {
		return u.PrimaryKey == true
	})

	if !usePrimaryKey || !existePrimaryKey {
		sq.errors = append(sq.errors, "Existen campos obligatorios sin information")
	}

	if len(sq.errors) > 0 {
		return nil, errors.New("No paso validaciones de filtrado")
	} else {
		return table_where, nil
	}
}

func (sp *Transactions) validaciones(item schema.Fields, value any) (any, error) {
	var val any
	var err error
	switch value.(type) {
	case string:
		val, err = caseString(value.(string), item.ValidateType.(schema.TypeStrings))
	case float64:
		val, err = caseFloat(value.(float64), item.ValidateType.(schema.TypeFloat64))
	case uint64:
		val, err = caseUint(value.(uint64), item.ValidateType.(schema.TypeUint64))
	case int64:
		val, err = caseInt(value.(int64), item.ValidateType.(schema.TypeInt64))
	case *time.Time:
		val = value
	default:
		val, err = nil, errors.New("tipo de dato no asignado")
	}
	return val, err
}

/*******************************Crud Multiples************************************/

func NewSqlExecMulti(database string) *SqlExecMultiples {
	return &SqlExecMultiples{
		database: database,
	}
}

func (sq *SqlExecMultiples) Set(s schema.Schema, datos ...schema.Models) *Transactions {
	key := len(sq.transaction)
	sq.transaction = append(sq.transaction, &Transactions{
		ob:     datos,
		schema: s,
	})

	return sq.transaction[key]
}

func (sq *SqlExecMultiples) Exec(params ...bool) error {
	cnn, err := new(internal.Connection).New(sq.database).NewPool()
	if err != nil {
		return err
	}
	defer cnn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Adquirir conexión del pool
	conn, err := cnn.Pool().Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error acquire pool: %w", err)
	}
	defer conn.Release()

	// Iniciar transacción
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error begin tx: %w", err)
	}

	//Manejar commit/rollback de forma segura
	success := false
	defer func() {
		if !success {
			_ = tx.Rollback(ctx)
		}
	}()

	cross := false
	if len(params) == 1 {
		cross = params[0]
	}

	// Ejecutar queries
	for _, t := range sq.transaction {
		if t.action == "" {
			tx.Rollback(ctx)
			return fmt.Errorf("existen datos sin procesar")
		}
		// time.Sleep(2 * time.Second)
		for _, item := range t.query {
			sqlPre := item["sqlPreparate"].(string)
			if cross && t.action == "UPDATE" {
				sqlPre = Query_Cross_Update(sqlPre)
			}
			valuesExec := item["valuesExec"].([]interface{})
			if _, err := tx.Exec(ctx, sqlPre, valuesExec...); err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("error sql %s: %s", t.action, logger.ManagerErrors{}.SqlCrud(err, t.schema.Name()).Error())
			}
		}
	}

	//Confirmar transacción
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commit: %w", err)
	}
	success = true
	return nil
}

func (sq *SqlExecMultiples) ExecTransaction(t *Transactions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sq.tx == nil {
		// Crear conexión
		cnn, err := new(internal.Connection).New(sq.database).NewPool()
		if err != nil {
			return err
		}
		defer cnn.Close()

		// Adquirir conexión del pool
		conn, err := cnn.Pool().Acquire(ctx)
		if err != nil {
			return fmt.Errorf("error acquire pool: %w", err)
		}
		defer conn.Release()

		// Iniciar transacción
		sq.tx, err = conn.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return fmt.Errorf("error begin tx: %w", err)
		}

	}

	// Manejo de rollback automático si algo falla
	success := false
	defer func() {
		if !success {
			_ = sq.tx.Rollback(ctx)
		}
	}()

	// Ejecutar todas las consultas del bloque
	for _, item := range t.query {
		sqlPre := item["sqlPreparate"].(string)
		valuesExec := item["valuesExec"].([]interface{})

		_, err := sq.tx.Exec(ctx, sqlPre, valuesExec...) // ← importante pasar ctx
		if err != nil {
			return fmt.Errorf("error SQL %s: %w", t.action, err)
		}
	}

	success = true
	return nil
}

func (sq *SqlExecMultiples) Commit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sq.tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commit transacción: %w", err)
	}
	return nil
}

/*
SetTransaction establece la información para nuevas transacciones, recibiendo directamente nuevas transacciones ya procesadas.

	Recibe uno o varias transacciones (s ...*Transaction) listas  para  ser ejecutadas.
	Retorna un array de  punteros a la transacción creadas.
	Parámetros
		* s {...*Transaction}: array de transacciones ya procesadas, listas para su ejecución
	Return
		- ([]*Transaction) retorna  []*Transaction
*/
func (sq *SqlExecMultiples) SetTransactions(s ...*Transactions) ([]*Transactions, error) {
	key := len(sq.transaction)
	var returned []*Transactions
	for _, v := range s {
		if v.action == "" {

			return nil, fmt.Errorf("existen datos sin procesar")
		}
		sq.transaction = append(sq.transaction, v)
		returned = append(returned, sq.transaction[key])
		key++
	}
	if len(returned) <= 0 {
		return nil, errors.New("se recibió datos sin ser procesados")
	}
	return returned, nil
}

/*
GetTransaction retorna las transacciones ya procesadas.

	Return
		- ([]*Transaction) retorna  []*Transaction
*/
func (sq *SqlExecMultiples) GetTransactions() []*Transactions {
	var returned []*Transactions
	for _, v := range sq.transaction {
		if v.action != "" {
			returned = append(returned, v)
		}
	}
	return returned
}
