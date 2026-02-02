package tables

import (
	"reflect"
	"strings"
	"time"

	"github.com/deybin/pgorm/schema"
	"github.com/google/uuid"
)

type ModelsSchema struct {
	table Models
}

type Models struct {
	AtCreate   *time.Time ` json:"atcreate" tag:"atcreate"  validate:"required" validateType:"" `
	Birthdate  *time.Time ` json:"birthdate" tag:"birthdate"  validate:"required" validateType:"" `
	Age        uint64     ` json:"age" tag:"age"  validate:"required" validateType:"max=100" `
	Amount     float64    ` json:"amount" tag:"amount"  validate:"required" validateType:"" `
	Credits    uint64     ` json:"credits" tag:"credits"  validate:"required;update;sum" validateType:"" `
	Passwords  string     ` json:"passwords" tag:"passwords"  validate:"required" validateType:"encrypt" `
	Key_Secret string     ` json:"key_secret" tag:"key_secret"  validate:"required" validateType:"cipher" `
	Document   string     ` json:"document" tag:"document"  validate:"required ;where" validateType:"min=1;max=11;case=lowercase" `
	Nombre     string     ` json:"nombre" tag:"nombre"  validate:"required;update" validateType:"case=lowercase" `
	Email      string     ` json:"email" tag:"correo electrónico"  validate:"required;update" validateType:"case=lowercase;expr=^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$" `
	Address    string     ` json:"address" tag:"address"  validate:"required" validateType:"case=lowercase" `
	Id         string     ` json:"id" tag:"id"  validate:"primaryKey; required;default" validateType:"" `
}

type ModelsFilter struct {
	Data       Models
	Conditions []schema.Where
}

func (s ModelsFilter) Name() string {
	t := reflect.TypeOf(s)
	return strings.ToLower(t.Name())
}

func (s Models) Name() string {
	t := reflect.TypeOf(s)
	return strings.ToLower(t.Name())
}

func (s ModelsSchema) Table() Models {
	return s.table
}

func (s ModelsSchema) Name() string {
	t := reflect.TypeOf(s.table)
	return strings.ToLower(t.Name())
}

func (s ModelsSchema) GetSchemaInsert() []schema.Fields {
	// s.table.Id = "f376f31c-c5ca-4e6f-a9ad-85e650bc061c"
	s.table.Id = uuid.NewString()
	// fmt.Println(s.table.Id)
	return schema.SchemaExe{}.GenerateSchema(s.table, schema.ActionInsert)
}

func (s ModelsSchema) GetSchemaUpdate() []schema.Fields {
	return schema.SchemaExe{}.GenerateSchema(s.table, schema.ActionUpdate)
}

func (s ModelsSchema) GetSchemaDelete() []schema.Fields {
	return schema.SchemaExe{}.GenerateSchema(s.table, schema.ActionDelete)
}
