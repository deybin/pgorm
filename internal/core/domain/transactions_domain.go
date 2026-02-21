package domain

import (
	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/migrator"
)

// type Actions string

// const (
// 	NONE   Actions = "NONE"
// 	INSERT Actions = "INSERT"
// 	UPDATE Actions = "UPDATE"
// 	DELETE Actions = "DELETE"
// )

type Transactions struct {
	ob       []migrator.Entity //datos para observación
	data     []map[string]any  //datos para insertar o actualizar o eliminar
	dataExec []adapters.DataExec
	schema   migrator.Schema
	action   adapters.Actions
	errors   []string
}

func NewTransaction(s migrator.Schema, datos ...migrator.Entity) Transactions {
	return Transactions{ob: datos, schema: s}
}

func (t Transactions) Schema() migrator.Schema {
	return t.schema
}

func (t Transactions) Action() adapters.Actions {
	return t.action
}

func (t Transactions) Data() []map[string]any {
	return t.data
}

func (t Transactions) Datos() []migrator.Entity {
	return t.ob
}

func (t Transactions) Query() []adapters.DataExec {
	return t.dataExec
}

func (t Transactions) Error() []string {
	return t.errors
}

func (t *Transactions) SetQuery(dataExec []adapters.DataExec) {
	t.dataExec = dataExec
}

func (t *Transactions) SetData(data []map[string]any) {

	t.data = data
}

func (t *Transactions) SetAction(action adapters.Actions) {
	t.action = action
}

/*******************************Crud Transactions************************************/
