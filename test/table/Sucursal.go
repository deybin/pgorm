package tables

import (
	"reflect"
	"strings"

	"github.com/deybin/pgorm"
)

type Schema_Sucursal struct {
	table Requ_Sucursal
}

type Requ_Sucursal struct {
	C_Sucu string ` json:"c_sucu" tag:"c_sucu"  validate:"required" validateType:"min=0,max=3,case=lowercase" `
	L_Sucu string ` json:"l_sucu" tag:"l_sucu"  validate:"required" validateType:"min=10,max=100,case=lowercase" `
	L_Dire string ` json:"l_dire" tag:"l_dire"  validate:"required" validateType:"min=20,max=200,case=lowercase" `
	C_Ubig string ` json:"c_ubig" tag:"c_ubig"  validate:"required" validateType:"min=1,max=6,case=lowercase" `
	N_Celu string ` json:"n_celu" tag:"n_celu"  validate:"required" validateType:"min=2,max=24,case=lowercase" `
	N_Tele string ` json:"n_tele" tag:"n_tele"  validate:"required" validateType:"min=2,max=24,case=lowercase" `
}

func (s *Schema_Sucursal) New() *Schema_Sucursal {
	return s
}

func (s *Schema_Sucursal) GetTableName() string {
	t := reflect.TypeOf(s.table)
	return strings.ToLower(t.Name())
}

func (s *Schema_Sucursal) GetSchemaInsert() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionInsert)
}

func (s *Schema_Sucursal) GetSchemaUpdate() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionUpdate)
}

func (s *Schema_Sucursal) GetSchemaDelete() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionDelete)
}
