package tables

import (
	"reflect"
	"strings"
	"time"

	"github.com/deybin/pgorm"
)

type Schema_Models struct {
	table Models
}

type Models struct {
	Atcreate   time.Time ` json:"atcreate" tag:"atcreate"  validate:"required" validateType:"" `
	Birthdate  time.Time ` json:"birthdate" tag:"birthdate"  validate:"required" validateType:"" `
	Age        uint64    ` json:"age" tag:"age"  validate:"required" validateType:"max=10" `
	Amount     float64   ` json:"amount" tag:"amount"  validate:"required" validateType:"" `
	Credits    uint64    ` json:"credits" tag:"credits"  validate:"required" validateType:"max=10" `
	Passwords  string    ` json:"passwords" tag:"passwords"  validate:"required" validateType:"case=lowercase" `
	Key_Secret string    ` json:"key_secret" tag:"key_secret"  validate:"required" validateType:"case=lowercase" `
	Document   string    ` json:"document" tag:"document"  validate:"required" validateType:"min=1,max=11,case=lowercase" `
	Nombre     string    ` json:"nombre" tag:"nombre"  validate:"required" validateType:"case=lowercase" `
	Address    string    ` json:"address" tag:"address"  validate:"required" validateType:"case=lowercase" `
	Id         string    ` json:"id" tag:"id"  validate:"required" validateType:"case=lowercase" `
}

func (s *Schema_Models) New() *Schema_Models {
	return s
}

func (s *Schema_Models) GetTableName() string {
	t := reflect.TypeOf(s.table)
	return strings.ToLower(t.Name())
}

func (s *Schema_Models) GetSchemaInsert() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionInsert)
}

func (s *Schema_Models) GetSchemaUpdate() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionUpdate)
}

func (s *Schema_Models) GetSchemaDelete() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionDelete)
}
