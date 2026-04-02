package test

import (
	"testing"

	"github.com/deybin/pgorm"
)

func TestModelsGenerateFile(t *testing.T) {

	pgorm.NewEntitySchema("master", "subscriptions")

	r := ""

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

func TestModelsGenerateFileWithSchema(t *testing.T) {

	pgorm.NewEntitySchemaWithSchemaName("platcont", "base_users", "other")

	r := ""

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}

// func TestCRUD_Model_Expr(t *testing.T) {

// 	dataInsert := map[string]interface{}{
// 		"document":   "719780401",
// 		"nombre":     "deybin yoni gil perez",
// 		"email":      "deybin.yoni@gmail.copm ",
// 		"address":    "av. general cordoba 427",
// 		"birthdate":  time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local),
// 		"age":        uint64(31),
// 		"amount":     15000.00,
// 		"credits":    int64(40),
// 		"passwords":  "Deybin04",
// 		"key_secret": "hola soy Deybin",
// 	}

// 	crud := pgorm.SqlExecSingle{}

// 	err := crud.New(&tables.ModelsSchema{}, dataInsert).Insert()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }
