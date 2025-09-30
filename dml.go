package pgorm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type SqlExecSingle struct {
	ob     []map[string]interface{} //datos para observación
	data   []map[string]interface{} //datos para insertar o actualizar o eliminar
	query  []map[string]interface{}
	schema Schema
	action string
}

type SqlExecMultiple struct {
	tx          pgx.Tx
	database    string
	transaction []*Transaction
}

type Transaction struct {
	ob     []map[string]interface{} //datos para observación
	data   []map[string]interface{} //datos para insertar o actualizar o eliminar
	query  []map[string]interface{}
	schema Schema
	action string
}

/*
New crea una nueva instancia de SqlExecSingle con el esquema y los datos proporcionados.

	Parámetros
		* s {Schema}: esquema de la tabla
		* datos {[]map[string]interface{}}: datos a insertar, actualizar o eliminar

	Return
		- (*SqlExecSingle) retorna  puntero *SqlExecSingle struct
*/
func (sq *SqlExecSingle) New(s Schema, datos ...map[string]interface{}) *SqlExecSingle {
	sq.ob = datos
	sq.schema = s
	return sq
}

/*
Valida los datos para insertar y crea el query para insertar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingle) Insert() error {
	sqlExec, data_insert, err := _insert(sq.schema.GetTableName(), sq.ob, sq.schema.GetSchemaInsert())
	if err != nil {
		return err
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
func (sq *SqlExecSingle) Update() error {
	sqlExec, data_update, err := _update(sq.schema.GetTableName(), sq.ob, sq.schema.GetSchemaUpdate())
	if err != nil {
		return err
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
func (sq *SqlExecSingle) Delete() error {
	sqlExec, data_delete, err := _delete(sq.schema.GetTableName(), sq.ob, sq.schema.GetSchemaDelete())
	if err != nil {
		return err
	}
	sq.query = sqlExec
	sq.data = data_delete
	sq.action = "DELETE"
	return nil
}

/*
Retorna los datos que se enviaron o enviaran para ser insertados, modificados o eliminados

	Return
		- []map[string]interface{}
*/
func (sq *SqlExecSingle) GetData() []map[string]interface{} {
	return sq.data
}

/*
Ejecuta el query

	Return
		- returns {error}: retorna errores ocurridos durante la ejecución
*/
func (sq *SqlExecSingle) Exec(database string, params ...bool) error {
	cnn, err := new(Connection).New(database).Pool()
	if err != nil {
		return err
	}

	cross := false
	if len(params) == 1 {
		cross = params[0]
	}
	dataExec := sq.query
	defer cnn.Close()
	for _, item := range dataExec {
		sqlPre := item["sqlPreparate"].(string)
		//fmt.Println("PREPARED: ", sqlPre)
		if cross {
			if sq.action == "UPDATE" {
				sqlPre = Query_Cross_Update(sqlPre)
			}
		}

		// fmt.Println("PREPARED: ", sqlPre)
		valuesExec := item["valuesExec"].([]interface{})

		if _, err_exec := cnn.pool.Exec(cnn.context, sqlPre, valuesExec...); err_exec != nil {
			return fmt.Errorf("error sql %s: %s", sq.action, err_exec.Error())
		}
	}
	return nil
}

/*
Crea una nueva instancia de SqlExecMultiple con el nombre de la base de datos proporcionado.

	Parámetros
	  * name {string}: database
	Return
	  - (*SqlExecMultiple) retorna  puntero *SqlExecMultiple struct
*/
func (sq *SqlExecMultiple) New(database string) *SqlExecMultiple {
	sq.database = database
	// sq.transaction = make(map[string]*Transaction)
	return sq
}

/*
SetInfo establece la información para una nueva transacción en SqlExecMultiple.

	Recibe un esquema (s Schema) y datos (datos ...map[string]interface{}) para la transacción.
	Retorna un puntero a la transacción creada.
	Parámetros
		* s {Schema}: esquema de la tabla
		* datos {[]map[string]interface{}}: datos a insertar, actualizar o eliminar
	Return
		- (*Transaction) retorna puntero *Transaction
*/
func (sq *SqlExecMultiple) SetInfo(s Schema, datos ...map[string]interface{}) *Transaction {
	key := len(sq.transaction)
	sq.transaction = append(sq.transaction, &Transaction{
		ob:     datos,
		schema: s,
	})

	return sq.transaction[key]
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
func (sq *SqlExecMultiple) SetTransaction(s ...*Transaction) ([]*Transaction, error) {
	key := len(sq.transaction)
	var returned []*Transaction
	for _, v := range s {
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
func (sq *SqlExecMultiple) GetTransactions() []*Transaction {
	var returned []*Transaction
	for _, v := range sq.transaction {
		if v.action != "" {
			returned = append(returned, v)
		}
	}
	return returned
}

/*
SetSqlSingle establece la información de un SqlExecSingle para crear una nueva Transaction, los datos del .SqlExecSingle ya deben estar procesados listo para la ejecución.

	Recibe SqlExecSingle (SqlExecSingle) listas  para  ser ejecutadas.
	Retorna un puntero a la transacción creada o bien el error si es que existiera
	Parámetros
		* s {...*Transaction}: array de transacciones ya procesadas, listas para su ejecución
	Return
		- (*Transaction,) retorna puntero *Transaction, si existe algún error
*/
func (sq *SqlExecMultiple) SetSqlSingle(s SqlExecSingle) (*Transaction, error) {
	key := len(sq.transaction)
	if s.action == "" {
		return nil, errors.New("datos de " + s.schema.GetTableName() + " aun no han sido procesados")
	}
	sq.transaction = append(sq.transaction, &Transaction{
		ob:     s.ob,
		data:   s.data,
		schema: s.schema,
		action: s.action,
		query:  s.query,
	})
	return sq.transaction[key], nil
}

/*
*
Ejecuta el query

	Return
		- (error): retorna errores ocurridos durante la ejecución
*/
func (sq *SqlExecMultiple) Exec(params ...bool) error {
	cnn, err := new(Connection).New(sq.database).Pool()
	if err != nil {
		return err
	}
	defer cnn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Adquirir conexión del pool
	conn, err := cnn.pool.Acquire(ctx)
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
		for _, item := range t.query {
			sqlPre := item["sqlPreparate"].(string)
			if cross && t.action == "UPDATE" {
				sqlPre = Query_Cross_Update(sqlPre)
			}
			valuesExec := item["valuesExec"].([]interface{})
			_, err := tx.Exec(ctx, sqlPre, valuesExec...)
			if err != nil {
				return fmt.Errorf("error SQL %s: %w", t.action, err)
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

func (t *Transaction) Insert() error {
	sqlExec, data_insert, err := _insert(t.schema.GetTableName(), t.ob, t.schema.GetSchemaInsert())
	if err != nil {
		return err
	}
	t.query = sqlExec
	t.data = data_insert
	t.action = "INSERT"
	return nil
}

func (t *Transaction) Update() error {
	sqlExec, data_update, err := _update(t.schema.GetTableName(), t.ob, t.schema.GetSchemaUpdate())
	if err != nil {
		return err
	}
	t.query = sqlExec
	t.data = data_update
	t.action = "UPDATE"
	return nil
}

func (t *Transaction) Delete() error {
	sqlExec, data_delete, err := _delete(t.schema.GetTableName(), t.ob, t.schema.GetSchemaDelete())
	if err != nil {
		return err
	}
	t.query = sqlExec
	t.data = data_delete
	t.action = "DELETE"
	return nil
}

func (t *Transaction) GetData() []map[string]interface{} {
	return t.data
}

func (sq *SqlExecMultiple) ExecTransaction(t *Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sq.tx == nil {
		// Crear conexión
		cnn, err := new(Connection).New(sq.database).Pool()
		if err != nil {
			return err
		}
		defer cnn.Close()

		// Adquirir conexión del pool
		conn, err := cnn.pool.Acquire(ctx)
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

func (sq *SqlExecMultiple) Commit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sq.tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commit transacción: %w", err)
	}
	return nil
}

func _insert(table string, data []map[string]interface{}, schema []Fields) ([]map[string]interface{}, []map[string]interface{}, error) {
	length := len(data)
	if length > 0 {
		var sqlExec = make([]map[string]interface{}, 0)
		var data_insert []map[string]interface{}

		for _, item := range data {
			preArray, err := _checkInsertSchema(schema, item)
			if err == nil {
				data_insert = append(data_insert, preArray)
				var column []string
				var values []string
				var i int
				var valuesExec []interface{}
				char := "$"
				for k, v := range preArray {
					i++
					column = append(column, k)
					values = append(values, fmt.Sprintf("%s%d", char, i))
					valuesExec = append(valuesExec, v)
				}

				sqlPreparate := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, strings.Join(column, ", "), strings.Join(values, ", "))
				sqlExec = append(sqlExec, map[string]interface{}{
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

func _update(table string, data []map[string]interface{}, schema []Fields) ([]map[string]interface{}, []map[string]interface{}, error) {
	length := len(data)

	if length > 0 {
		var sqlExec = make([]map[string]interface{}, 0)
		var data_update []map[string]interface{}
		for _, item := range data {
			where := make(map[string]interface{})

			length_where := 0
			if item["where"] != nil {
				where = item["where"].(map[string]interface{})
				length_where = len(where)
				delete(item, "where")
			}
			preArray, err := _checkUpdate(schema, item)
			if err != nil {
				return nil, nil, err
			}

			if len(preArray) <= 0 {
				continue
			}

			preArray_where := make(map[string]interface{})
			if length_where > 0 {
				preArray, err := _checkWhere(schema, where)
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

			if length_where > 0 {
				length_newMapWhere := len(preArray_where)
				var wheres []string
				for k, v := range preArray_where {
					i++
					wheres = append(wheres, k+" = "+char+strconv.FormatUint(i, 10))
					valuesExec = append(valuesExec, v)
				}
				if length_newMapWhere > 0 {
					sqlWherePreparateUpdate = "WHERE " + strings.Join(wheres, " AND ")
				}
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

func _delete(table string, data []map[string]interface{}, schema []Fields) ([]map[string]interface{}, []map[string]interface{}, error) {
	length := len(data)

	if length > 0 {
		var sqlExec = make([]map[string]interface{}, 0)
		var data_delete []map[string]interface{}
		for _, item := range data {

			preArray, err := _checkWhere(schema, item)
			if err != nil {
				return nil, nil, err
			}

			data_delete = append(data_delete, preArray)
			var lineSqlExec = make(map[string]interface{}, 2)
			sqlWherePreparateDelete := ""
			var i int
			var p uint64
			length_newMap := len(preArray)
			var valuesExec []interface{}
			if length_newMap > 0 {
				sqlWherePreparateDelete += " WHERE "
			}
			char := "$"
			for k, v := range preArray {
				p++
				if i+1 < length_newMap {
					// sqlWherePreparateUpdate += fmt.Sprintf("%s = '%s' AND ", ke, va)
					sqlWherePreparateDelete += k + " = " + char + strconv.FormatUint(p, 10) + " AND "
				} else {
					//sqlWherePreparateUpdate += fmt.Sprintf("%s = '%s'", ke, va)
					sqlWherePreparateDelete += k + " = " + char + strconv.FormatUint(p, 10)
				}
				valuesExec = append(valuesExec, v)
				i++
			}

			sqlPreparate := fmt.Sprintf("DELETE FROM %s %s", table, sqlWherePreparateDelete)
			lineSqlExec["sqlPreparate"] = sqlPreparate
			lineSqlExec["valuesExec"] = valuesExec
			sqlExec = append(sqlExec, lineSqlExec)
		}
		return sqlExec, data_delete, nil
	} else {
		return nil, nil, errors.New("no existen datos para actualizar")
	}
}

func _checkInsertSchema(schema []Fields, tabla_map map[string]interface{}) (map[string]interface{}, error) {

	// var err_cont uint64 = 0
	var err_cont uint
	var error string

	data := make(map[string]interface{})

	for _, item := range schema {
		isNil := tabla_map[item.Name] == nil
		defaultIsNil := item.Default == nil
		if !isNil {
			value := tabla_map[item.Name]
			new_value, err := strconvDataType(string(item.Type), value)
			if err != nil {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s %s", err_cont, item.Description, err.Error())
			}
			var val interface{}
			val, err = validaciones(item, new_value)
			if err == nil {
				data[item.Name] = val
			} else {
				err_cont++
				error += fmt.Sprintf("%d.- Se encontró fallas al validar el campo %s \n %s\n", err_cont, item.Description, err.Error())
			}
		} else {
			if !defaultIsNil {
				data[item.Name] = item.Default
			} else {
				if item.Required {
					err_cont++
					error += fmt.Sprintf("%d.- El campo %s es Requerido\n", err_cont, item.Description)
				}
			}
		}

	}
	if err_cont > 0 {
		return nil, errors.New(error)
	} else {
		return data, nil
	}
}

func _checkUpdate(schema []Fields, tabla_map map[string]interface{}) (map[string]interface{}, error) {
	var err_cont uint
	var error string
	data := make(map[string]interface{})
	for _, item := range schema {
		isNil := tabla_map[item.Name] == nil
		if !isNil {
			if item.Update {
				value := tabla_map[item.Name]
				new_value, err := strconvDataType(string(item.Type), value)
				if err != nil {
					err_cont++
					error += fmt.Sprintf("%d.- El campo %s %s", err_cont, item.Description, err.Error())
				}
				if item.Type == "string" {
					if new_value.(string) == "" {
						if !item.Empty {
							err_cont++
							error += fmt.Sprintf("%d.- El campo %s no puede estar vació\n", err_cont, item.Description)
						}
					}
				}

				var val interface{}
				val, err = validaciones(item, new_value)

				if err == nil {
					keyName := item.Name

					switch item.ArithmeticOperations {
					case Sum:
						keyName = fmt.Sprintf("ADD_%s_SUMA", item.Name)
					case Subtraction:
						keyName = fmt.Sprintf("ADD_%s_SUBTRACTION", item.Name)
					case Multiply:
						keyName = fmt.Sprintf("ADD_%s_MULTIPLY", item.Name)
					case Divide:
						keyName = fmt.Sprintf("ADD_%s_DIVIDE", item.Name)
					}

					data[keyName] = val
				} else {
					err_cont++
					error += fmt.Sprintf("%d.- Se encontró fallas al validar el campo %s \n %s\n", err_cont, item.Description, err.Error())
				}
			} else {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s no puede ser modificado\n", err_cont, item.Description)
			}
		}
	}
	if err_cont > 0 {
		return nil, errors.New(error)
	} else {
		return data, nil
	}
}

func _checkWhere(schema []Fields, table_where map[string]interface{}) (map[string]interface{}, error) {
	var err_cont uint
	var error string
	data := make(map[string]interface{})
	for _, item := range schema {
		isNil := table_where[item.Name] == nil
		if !isNil {
			value := table_where[item.Name]
			if !item.Where && !item.PrimaryKey {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s no puede ser utilizado de esta forma\n", err_cont, item.Description)
			} else {
				if value.(string) == "" {
					err_cont++
					error += fmt.Sprintf("%d.- El campo %s esta vació verificar\n", err_cont, item.Description)
				} else {
					data[item.Name] = value
				}

			}
		} else {
			if item.PrimaryKey {
				err_cont++
				error += fmt.Sprintf("%d.- El campo %s es obligatorio\n", err_cont, item.Description)
			}
		}
	}
	if err_cont > 0 {
		return nil, errors.New(error)
	} else {
		return data, nil
	}
}

func caseString(value string, schema TypeStrings) (string, error) {
	value = strings.TrimSpace(value)
	if schema.Expr != nil {
		if !schema.Expr.MatchString(value) {
			return "", errors.New("no cumple con las características")
		}
	}

	if schema.Encriptar {
		result, _ := bcrypt.GenerateFromPassword([]byte(value), 13)
		value = string(result)
		return value, nil
	}

	if schema.Cifrar {
		hash, _ := AesEncrypt_PHP([]byte(value), GetKey_PrivateCrypto())
		value = hash
		return value, nil
	}

	if schema.Min > 0 {
		if len(value) < schema.Min {
			return "", fmt.Errorf("no Cumple los caracteres mínimos que debe tener (%v)", schema.Min)
		}
	}

	if schema.Max > 0 {
		if len(value) > schema.Max {
			return "", fmt.Errorf("no Cumple los caracteres máximos que debe tener (%v)", schema.Max)
		}
	}

	if schema.UpperCase {
		value = strings.ToUpper(value)
	} else if schema.LowerCase {
		value = strings.ToLower(value)
	}
	return value, nil
}

func caseDate(value time.Time, _ TypeDate) (time.Time, error) {
	return value, nil
}

func caseFloat(value float64, schema TypeFloat64) (float64, error) {
	error := ""
	err_cont := 0
	if schema.Menor != 0 {
		if value <= schema.Menor {
			err_cont++
			error += fmt.Sprintf("- No puede se menor a %f", schema.Menor)
		}
	}
	if schema.Mayor != 0 {
		if value >= schema.Mayor {
			err_cont++
			error += fmt.Sprintf("- No puede se mayor a %f", schema.Mayor)
		}
	}
	if !schema.Negativo {
		if value < float64(0) {
			err_cont++
			error += "- No puede ser negativo"
		}
	}
	if schema.Porcentaje {
		value = value / float64(100)
	}
	if err_cont > 0 {
		return 0, errors.New(error)
	} else {
		return value, nil
	}
}

func caseInt(value int64, schema TypeInt64) (int64, error) {
	error := ""
	err_cont := 0
	if !schema.Negativo {
		if value < int64(0) {
			err_cont++
			error += "- No puede ser negativo"
		}
	}
	if schema.Min != 0 {
		if value < schema.Min {
			err_cont++
			error += fmt.Sprintf("- No puede se menor a %d", schema.Min)
		}
	}
	if schema.Max != 0 {
		if value > schema.Max {
			err_cont++
			error += fmt.Sprintf("- No puede se mayor a %d", schema.Max)
		}
	}
	if err_cont > 0 {
		return int64(0), errors.New(error)
	} else {
		return value, nil
	}
}

func caseUint(value uint64, schema TypeUint64) (uint64, error) {
	if schema.Max > 0 {
		if value > schema.Max {
			return 0, errors.New("- no esta en el rango permitido")
		}
	}
	return value, nil
}

func strconvDataType(targetType string, value interface{}) (interface{}, error) {
	//fmt.Println(reflect.TypeOf(value).String(), targetType)
	switch targetType {
	case "string":
		switch v := value.(type) {
		case string:
			return v, nil
		case fmt.Stringer:
			return v.String(), nil
		default:
			return fmt.Sprintf("%v", value), nil
		}

	case "float64":
		switch v := value.(type) {
		case float64:
			return v, nil
		case int64:
			return float64(v), nil
		case uint64:
			return float64(v), nil
		case string:
			return strconv.ParseFloat(v, 64)
		default:
			return nil, fmt.Errorf("no se puede convertir %T a float64", value)
		}

	case "int64":
		switch v := value.(type) {
		case int64:
			return v, nil
		case int:
			return int64(v), nil
		case float64:
			return int64(v), nil
		case uint64:
			return int64(v), nil
		case string:
			return strconv.ParseInt(v, 10, 64)
		default:
			return nil, fmt.Errorf("no se puede convertir %T a int64", value)
		}

	case "uint64":
		switch v := value.(type) {
		case uint64:
			return v, nil
		case float64:
			return uint64(v), nil
		case int:
			return uint64(v), nil
		case int64:
			return uint64(v), nil
		case string:
			return strconv.ParseUint(v, 10, 64)
		default:
			return nil, fmt.Errorf("no se puede convertir %T a uint64", value)
		}

	case "time.Time":
		switch v := value.(type) {
		case time.Time:
			return v, nil
		case string:
			// Puedes cambiar el formato si usas otro
			layout := "02/01/2006"
			return time.Parse(layout, v)
		default:
			return nil, fmt.Errorf("no se puede convertir %T a time.Time", value)
		}

	default:
		return nil, fmt.Errorf("tipo de destino no soportado: %s", targetType)
	}
}

func validaciones(item Fields, new_value interface{}) (interface{}, error) {
	var val interface{}
	var err error
	switch item.Type {
	case "string":
		val, err = caseString(new_value.(string), item.ValidateType.(TypeStrings))
	case "float64":
		val, err = caseFloat(new_value.(float64), item.ValidateType.(TypeFloat64))
	case "uint64":
		val, err = caseUint(new_value.(uint64), item.ValidateType.(TypeUint64))
	case "int64":
		val, err = caseInt(new_value.(int64), item.ValidateType.(TypeInt64))
	case "time.Time":
		val, err = caseDate(new_value.(time.Time), item.ValidateType.(TypeDate))
	default:
		val, err = nil, errors.New("tipo de dato no asignado")
	}
	return val, err
}
