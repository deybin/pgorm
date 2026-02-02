package test

import (
	"os"
	"testing"
	"time"

	"github.com/deybin/pgorm"
	"github.com/deybin/pgorm/schema"
	tables "github.com/deybin/pgorm/test/table"
	"github.com/joho/godotenv"
)

// inserta,actualiza y elimina datos solo de una tabla
func TestCRUD_IUD_Single(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")
	birth := time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local)
	dataInsert := tables.Models{
		Id:         "369d48ab-d881-4e22-9f14-8e9af67da9aa",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        uint64(31),
		Amount:     17000.00,
		Credits:    uint64(40),
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_Secret: "hola soy Deybin",
	}

	// println(dataInsert.GetSchemaInsert())
	crudInsert := pgorm.NewSqlExecSingles(&tables.ModelsSchema{}, dataInsert)

	err := crudInsert.Transaction().Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crudInsert.Exec(database)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	time.Sleep(10 * time.Second)

	dataUpdate := tables.ModelsFilter{
		Data: tables.Models{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []schema.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	crudUpdate := pgorm.NewSqlExecSingles(&tables.ModelsSchema{}, dataUpdate)

	if err := crudUpdate.Transaction().Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := crudUpdate.Exec(database); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	time.Sleep(10 * time.Second)

	dataDelete := []schema.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	crudDelete := pgorm.NewSqlExecSingles(&tables.ModelsSchema{})

	if err := crudDelete.Transaction().Delete(dataDelete...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := crudDelete.Exec(database); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
}

func TestCRUD_IUD_Multiple(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")
	birth := time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local)
	dataInsert := tables.Models{
		Id:         "550e8400-e29b-41d4-a716-446655440011",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        uint64(31),
		Amount:     15000.00,
		Credits:    uint64(40),
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_Secret: "hola soy Deybin",
	}

	dataInsert2 := tables.Models2{
		Id:         "550e8400-e29b-41d4-a716-446655440012",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        uint64(31),
		Amount:     15000.00,
		Credits:    uint64(40),
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_Secret: "hola soy Deybin",
	}

	// println(dataInsert.GetSchemaInsert())
	crud := pgorm.NewSqlExecMulti(database)

	table := crud.Set(&tables.ModelsSchema{}, dataInsert)

	if err := table.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	table2 := crud.Set(&tables.Models2Schema{}, dataInsert2)
	if err := table2.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	dataUpdate2 := tables.Models2Filter{
		Data: tables.Models2{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []schema.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "550e8400-e29b-41d4-a716-446655440012"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	table3 := crud.Set(&tables.Models2Schema{}, dataUpdate2)
	if err := table3.Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataUpdate := tables.ModelsFilter{
		Data: tables.Models{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []schema.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "550e8400-e29b-41d4-a716-446655440011"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	crudUpdate := pgorm.NewSqlExecSingles(&tables.ModelsSchema{}, dataUpdate)
	if err := crudUpdate.Transaction().Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	crud.SetTransactions(crudUpdate.Transaction())

	dataDelete1 := []schema.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	table4 := crud.Set(&tables.Models2Schema{})
	if err := table4.Delete(dataDelete1...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataDelete := []schema.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	crudDelete := pgorm.NewSqlExecSingles(&tables.ModelsSchema{})

	if err := crudDelete.Transaction().Delete(dataDelete...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	crud.SetTransactions(crudDelete.Transaction())

	if err := crud.Exec(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

// func TestCRUD_Insert_SingleDefault(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")

// 	dataInsert := map[string]interface{}{
// 		"document":   "719780401",
// 		"nombre":     "deybin yoni gil perez",
// 		"address":    "av. general cordoba 427",
// 		"birthdate":  time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local),
// 		"age":        uint64(31),
// 		"amount":     15000.00,
// 		"credits":    int64(40),
// 		"passwords":  "Deybin04",
// 		"key_secret": "hola soy Deybin",
// 	}

// 	crud := pgorm.SqlExecSingle{}

// 	err := crud.New(new(tables.Models).New(), dataInsert).Insert()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Exec(database)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// }

// // insertar datos con valores DIFERENTES a los solicitados por el modelo
// func TestCRUD_Insert_SingleType(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")
// 	dataInsert := map[string]interface{}{
// 		"id":         "412904da-43ce-4aab-b79b-beea21c5d86e",
// 		"document":   "719780401",
// 		"nombre":     "deybin yoni gil perez",
// 		"address":    "av. general cordoba",
// 		"birthdate":  "04/04/1994",
// 		"age":        31,
// 		"amount":     "15000.00",
// 		"credits":    40,
// 		"passwords":  "Deybin04",
// 		"key_secret": "hola soy Deybin",
// 	}
// 	crud := pgorm.SqlExecSingle{}

// 	err := crud.New(new(tables.Models).New(), dataInsert).Insert()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Exec(database)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// }

// // actualizar campos que te permite la actualización según el modelo
// func TestCRUD_Update_Single_Update(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")

// 	dataUpdate := map[string]interface{}{
// 		"nombre":  "Nuevo Nombre Actualizado",
// 		"credits": 1,
// 		"where":   map[string]interface{}{"id": "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
// 	}

// 	crud := pgorm.SqlExecSingle{}
// 	err := crud.New(new(tables.Models).New(), dataUpdate).Update()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Exec(database, true)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }

// // actualizar campos que no se permite la actualización según el modelo
// func TestCRUD_SingleError_Update(t *testing.T) {
// 	godotenv.Load()
// 	//database := os.Getenv("ENV_TES_DB")

// 	dataUpdate := map[string]interface{}{
// 		"document": "719780401",
// 		"where":    map[string]interface{}{"id": "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
// 	}

// 	crud := pgorm.SqlExecSingle{}
// 	err := crud.New(new(tables.Models).New(), dataUpdate).Update()
// 	if err.Error() != "al realizar validaciones se filtro datos y se quedo sin información para actualizar" {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 	}
// }

// // eliminar registro filtrando por primary key
// func TestCRUD_Single_Delete(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")

// 	dataDelete := map[string]interface{}{
// 		"id": "412904da-43ce-4aab-b79b-beea21c5d86e",
// 	}

// 	crud := pgorm.SqlExecSingle{}
// 	err := crud.New(new(tables.Models).New(), dataDelete).Delete()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Exec(database)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }

// // eliminar registro filtrando por primary key y un campo que permita filtrado
// func TestCRUD_SingleError_Delete(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")

// 	dataDelete := map[string]interface{}{
// 		"id":     "39ecfb80-3f46-4a1c-a019-d275cefd17ba",
// 		"nombre": "39ecfb80-3f46-4a1c-a019-d275cefd17ba",
// 	}

// 	crud := pgorm.SqlExecSingle{}
// 	err := crud.New(new(tables.Models).New(), dataDelete).Delete()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Exec(database)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }

// func TestCRUD_Multiple(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")
// 	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
// 		"id":         "412904da-43ce-4aab-b79b-beea21c5d86a",
// 		"document":   "11111111",
// 		"nombre":     "deybin yoni gil perez",
// 		"address":    "av. general cordoba",
// 		"birthdate":  "04/04/1994",
// 		"age":        31,
// 		"amount":     "15000.00",
// 		"credits":    40,
// 		"passwords":  "Deybin04",
// 		"key_secret": "hola soy Deybin",
// 	})
// 	dataInsert = append(dataInsert, map[string]interface{}{
// 		"id":         "412904da-43ce-4aab-b79b-beea21c5d86e",
// 		"document":   "22222222",
// 		"nombre":     "deybin yoni gil perez",
// 		"address":    "av. general cordoba",
// 		"birthdate":  "04/04/1994",
// 		"age":        31,
// 		"amount":     "15000.00",
// 		"credits":    40,
// 		"passwords":  "Deybin04",
// 		"key_secret": "hola soy Deybin",
// 	})

// 	dataInsert2 := map[string]interface{}{
// 		"id":         "412904da-43ce-4aab-b79b-beea21c5d87a",
// 		"document":   "33333333",
// 		"nombre":     "deybin yoni gil perez",
// 		"address":    "av. general cordoba",
// 		"birthdate":  "04/04/1994",
// 		"age":        31,
// 		"amount":     "15000.00",
// 		"credits":    40,
// 		"passwords":  "Deybin04",
// 		"key_secret": "hola soy Deybin",
// 	}

// 	crud := pgorm.SqlExecMultiple{}
// 	crud.New(database)

// 	trSucursal := crud.SetInfo(new(tables.Models).New(), dataInsert...)
// 	trAlmacen := crud.SetInfo(new(tables.Models2).New(), dataInsert2)

// 	if err := trAlmacen.Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	if err := trSucursal.Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	if err := crud.Exec(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }

// func TestCRUD_Multiple_transaction(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")
// 	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "004",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "indefinido",
// 	})
// 	dataInsert = append(dataInsert, map[string]interface{}{
// 		"c_sucu": "005",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "indefinido",
// 	})
// 	dataInsert = append(dataInsert, map[string]interface{}{
// 		"c_sucu": "003",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "indefinido",
// 	})

// 	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "004",
// 		"c_alma": "001",
// 		"l_alma": "sucursal de prueba",
// 	})

// 	// schema, tableName := table.GetSucursal()
// 	// schemaAlma, tableNameAlma := table.GetAlmacen()
// 	crud := new(pgorm.SqlExecMultiple).New(database)
// 	TransactionAlmacen := crud.New(database).SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma...)
// 	TransactionSucursal := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert...)

// 	err := TransactionAlmacen.Insert()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// 	err = TransactionSucursal.Insert()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.ExecTransaction(TransactionAlmacen)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.ExecTransaction(TransactionSucursal)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Commit()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// }

// func TestCRUD_Multiple_Set_Transaction(t *testing.T) {

// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")

// 	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "003",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "sin información",
// 	})

// 	dataInsert = append(dataInsert, map[string]interface{}{
// 		"c_sucu": "004",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "sin información",
// 	})

// 	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "003",
// 		"c_alma": "001",
// 		"l_alma": "sucursal de prueba",
// 	})
// 	dataInsertAlma = append(dataInsertAlma, map[string]interface{}{
// 		"c_sucu": "004",
// 		"c_alma": "001",
// 		"l_alma": "sucursal de prueba",
// 	})

// 	dataInsert2 := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "005",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "sin información",
// 	})

// 	dataInsertAlma2 := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "004",
// 		"c_alma": "001",
// 		"l_alma": "sucursal de prueba",
// 	})

// 	crud := pgorm.SqlExecMultiple{}
// 	crud.New(database)

// 	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// 	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	crud2 := pgorm.SqlExecMultiple{}
// 	crud2.New(database)

// 	if err := crud2.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert2...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// 	if err := crud2.SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma2...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	if _, err := crud.SetTransaction(crud2.GetTransactions()...); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	if err := crud.Exec(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }

// func TestCRUD_Multiple_Set_Single(t *testing.T) {

// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")

// 	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "003",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "sin información",
// 	})

// 	dataInsert = append(dataInsert, map[string]interface{}{
// 		"c_sucu": "004",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "sin información",
// 	})

// 	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "003",
// 		"c_alma": "001",
// 		"l_alma": "sucursal de prueba",
// 	})
// 	dataInsertAlma = append(dataInsertAlma, map[string]interface{}{
// 		"c_sucu": "004",
// 		"c_alma": "001",
// 		"l_alma": "sucursal de prueba",
// 	})

// 	dataInsert2 := append([]map[string]interface{}{}, map[string]interface{}{
// 		"c_sucu": "005",
// 		"l_sucu": "sucursal de prueba",
// 		"l_dire": "sin información",
// 	})

// 	crud := pgorm.SqlExecMultiple{}
// 	crud.New(database)

// 	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}
// 	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	crud2 := pgorm.SqlExecSingle{}
// 	crud2.New(new(tables.Schema_Sucursal).New(), dataInsert2...)

// 	if err := crud2.New(new(tables.Schema_Sucursal).New(), dataInsert2...).Insert(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	if _, err := crud.SetSqlSingle(crud2); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	if err := crud.Exec(); err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }

// func TestCRUD_Single_ADD_SUM(t *testing.T) {
// 	godotenv.Load()
// 	database := os.Getenv("ENV_TES_DB")
// 	dataUpdate := map[string]interface{}{
// 		"credits": uint64(1),
// 		"where":   map[string]interface{}{"id": "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
// 	}

// 	crud := pgorm.SqlExecSingle{}
// 	err := crud.New(new(tables.Models).New(), dataUpdate).Update()
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// 	err = crud.Exec(database, true)
// 	if err != nil {
// 		t.Errorf("se esperaba este error: %s", err.Error())
// 		return
// 	}

// }
