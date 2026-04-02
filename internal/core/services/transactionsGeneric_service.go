package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/builder"
	"github.com/deybin/pgorm/internal/core/domain"
	"github.com/deybin/pgorm/migrator"
)

type SqlExecSingles struct {
	Transactions domain.Transactions
	errors       []string
}

type SqlExecMultiples struct {
	transaction []*SqlExecSingles
}

// /*
// New crea una nueva instancia de SqlExecSingle con el esquema y los datos proporcionados.

// 	Parámetros
// 		* s {Schema}: esquema de la tabla
// 		* datos {[]map[string]interface{}}: datos a insertar, actualizar o eliminar

// 	Return
// 		- (*SqlExecSingle) retorna  puntero *SqlExecSingle struct
// */
// func NewSqlExecSingles(s migrator.Schema, datos ...migrator.Entity) *SqlExecSingles {
// 	return &SqlExecSingles{Transactions: domain.NewTransaction(s, datos...)}
// }

func (sq *SqlExecSingles) Transaction() *domain.Transactions {
	return &sq.Transactions
}

func (sq *SqlExecSingles) Schema() migrator.Schema {
	return sq.Transactions.Schema()
}

func (sq *SqlExecSingles) Datos() []migrator.Entity {
	return sq.Transactions.Datos()
}

/*
Valida los datos para insertar y crea el query para insertar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingles) Insert() error {
	sq.Transaction().SetAction(adapters.INSERT)
	if err := builder.BuilderInsertGeneric(sq.Transaction()); err != nil {
		return err
	}
	return nil
}

/*
Valida los datos para actualizar y crea el query para actualizar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingles) Update() error {
	sq.Transaction().SetAction(adapters.UPDATE)
	if err := builder.BuilderUpdateGeneric(sq.Transaction()); err != nil {
		return err
	}
	return nil
}

/*
Valida los datos para Eliminar y crea el query para Eliminar

	Return
		- (error): retorna errores ocurridos en la validación
*/
func (sq *SqlExecSingles) Delete(dataDelete ...migrator.Where) error {
	sq.Transaction().SetAction(adapters.DELETE)
	if err := builder.BuilderDeleteGeneric(sq.Transaction(), dataDelete); err != nil {
		return errors.New(strings.Join(sq.errors, "; "))
	}
	return nil
}

/*******************************Crud Multiples************************************/

func NewSqlExecMulti() *SqlExecMultiples {
	return &SqlExecMultiples{}
}

func (sq *SqlExecMultiples) Set(ses *SqlExecSingles) *SqlExecSingles {
	key := len(sq.transaction)
	sq.transaction = append(sq.transaction, ses)
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
func (sq *SqlExecMultiples) SetTransactions(s ...*SqlExecSingles) ([]*SqlExecSingles, error) {
	key := len(sq.transaction)
	var returned []*SqlExecSingles
	for _, v := range s {
		if v.Transaction().Action() == adapters.NONE {

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
		- ([]*SqlExecSingles) retorna  []*SqlExecSingles
*/
func (sq *SqlExecMultiples) GetTransactions() []*SqlExecSingles {
	var returned []*SqlExecSingles
	for _, v := range sq.transaction {
		if v.Transaction().Action() != adapters.NONE {
			returned = append(returned, v)
		}
	}
	return returned
}

/*
GetTransaction retorna las transacciones ya procesadas.

	Return
		- ([]*SqlExecSingles) retorna  []*SqlExecSingles
*/
func (sq *SqlExecMultiples) DataExec() [][]adapters.DataExec {
	var returned [][]adapters.DataExec
	for _, v := range sq.transaction {
		if v.Transaction().Action() != adapters.NONE {
			returned = append(returned, v.Transaction().Query())
		}
	}
	return returned
}
