package tables

import (
	"reflect"
	"strings"
	"time"

	"github.com/deybin/pgorm"
)

type Schema_Models2 struct {
	table Models2
}

type Models2 struct {
	Atcreate   time.Time ` json:"atcreate" tag:"atcreate"  validate:"required" validateType:"" `
	Birthdate  time.Time ` json:"birthdate" tag:"birthdate"  validate:"required" validateType:"" `
	Age        uint64    ` json:"age" tag:"age"  validate:"required" validateType:"max=100" `
	Amount     float64   ` json:"amount" tag:"amount"  validate:"required" validateType:"" `
	Credits    uint64    ` json:"credits" tag:"credits"  validate:"required,update,sum" validateType:"" `
	Passwords  string    ` json:"passwords" tag:"passwords"  validate:"required" validateType:"encrypt" `
	Key_Secret string    ` json:"key_secret" tag:"key_secret"  validate:"required" validateType:"cipher" `
	Document   string    ` json:"document" tag:"document"  validate:"required" validateType:"min=1,max=11,case=lowercase" `
	Nombre     string    ` json:"nombre" tag:"nombre"  validate:"required,update,where" validateType:"case=lowercase" `
	Address    string    ` json:"address" tag:"address"  validate:"required" validateType:"case=lowercase" `
	Id         string    ` json:"id" tag:"id"  validate:"primaryKey" validateType:"" `
}

func (s *Schema_Models2) New() *Schema_Models2 {
	return s
}

func (s *Schema_Models2) GetTableName() string {
	t := reflect.TypeOf(s.table)
	return strings.ToLower(t.Name())
}

func (s *Schema_Models2) GetSchemaInsert() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionInsert)
}

func (s *Schema_Models2) GetSchemaUpdate() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionUpdate)
}

func (s *Schema_Models2) GetSchemaDelete() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionDelete)
}
