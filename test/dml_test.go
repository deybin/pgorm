package test

import (
	"os"
	"testing"
	"time"

	"github.com/deybin/pgorm"
	tables "github.com/deybin/pgorm/test/table"
	"github.com/joho/godotenv"
)

// insertar datos con valores IGUALES a los solicitados por el modelo
func TestCRUD_Insert_Single(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")

	dataInsert := map[string]interface{}{
		"id":         "369d48ab-d881-4e22-9f14-8e9af67da9aa",
		"document":   "719780401",
		"nombre":     "deybin yoni gil perez",
		"address":    "av. general cordoba 427",
		"birthdate":  time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local),
		"age":        uint64(31),
		"amount":     15000.00,
		"credits":    int64(40),
		"passwords":  "Deybin04",
		"key_secret": "hola soy Deybin",
	}

	crud := pgorm.SqlExecSingle{}

	err := crud.New(new(tables.Schema_Models).New(), dataInsert).Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec(database)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
}

// insertar datos con valores DIFERENTES a los solicitados por el modelo
func TestCRUD_Insert_SingleType(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")
	dataInsert := map[string]interface{}{
		"id":         "412904da-43ce-4aab-b79b-beea21c5d86e",
		"document":   "719780401",
		"nombre":     "deybin yoni gil perez",
		"address":    "av. general cordoba",
		"birthdate":  "04/04/1994",
		"age":        31,
		"amount":     "15000.00",
		"credits":    40,
		"passwords":  "Deybin04",
		"key_secret": "hola soy Deybin",
	}
	crud := pgorm.SqlExecSingle{}

	err := crud.New(new(tables.Schema_Models).New(), dataInsert).Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec(database)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
}

// actualizar campos que te permite la actualización según el modelo
func TestCRUD_Update_Single_Update(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")

	dataUpdate := map[string]interface{}{
		"nombre":  "Nuevo Nombre Actualizado",
		"credits": 1,
		"where":   map[string]interface{}{"id": "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
	}

	crud := pgorm.SqlExecSingle{}
	err := crud.New(new(tables.Schema_Models).New(), dataUpdate).Update()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec(database, true)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

// actualizar campos que no se permite la actualización según el modelo
func TestCRUD_SingleError_Update(t *testing.T) {
	godotenv.Load()
	//database := os.Getenv("ENV_TES_DB")

	dataUpdate := map[string]interface{}{
		"document": "719780401",
		"where":    map[string]interface{}{"id": "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
	}

	crud := pgorm.SqlExecSingle{}
	err := crud.New(new(tables.Schema_Models).New(), dataUpdate).Update()
	if err.Error() != "al realizar validaciones se filtro datos y se quedo sin información para actualizar" {
		t.Errorf("se esperaba este error: %s", err.Error())
	}
}

// eliminar registro filtrando por primary key
func TestCRUD_Single_Delete(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")

	dataDelete := map[string]interface{}{
		"id": "412904da-43ce-4aab-b79b-beea21c5d86e",
	}

	crud := pgorm.SqlExecSingle{}
	err := crud.New(new(tables.Schema_Models).New(), dataDelete).Delete()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec(database)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

// eliminar registro filtrando por primary key y un campo que permita filtrado
func TestCRUD_SingleError_Delete(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")

	dataDelete := map[string]interface{}{
		"id":     "39ecfb80-3f46-4a1c-a019-d275cefd17ba",
		"nombre": "39ecfb80-3f46-4a1c-a019-d275cefd17ba",
	}

	crud := pgorm.SqlExecSingle{}
	err := crud.New(new(tables.Schema_Models).New(), dataDelete).Delete()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec(database)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Multiple(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")
	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
		"id":         "412904da-43ce-4aab-b79b-beea21c5d86a",
		"document":   "11111111",
		"nombre":     "deybin yoni gil perez",
		"address":    "av. general cordoba",
		"birthdate":  "04/04/1994",
		"age":        31,
		"amount":     "15000.00",
		"credits":    40,
		"passwords":  "Deybin04",
		"key_secret": "hola soy Deybin",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"id":         "412904da-43ce-4aab-b79b-beea21c5d86e",
		"document":   "22222222",
		"nombre":     "deybin yoni gil perez",
		"address":    "av. general cordoba",
		"birthdate":  "04/04/1994",
		"age":        31,
		"amount":     "15000.00",
		"credits":    40,
		"passwords":  "Deybin04",
		"key_secret": "hola soy Deybin",
	})

	dataInsert2 := map[string]interface{}{
		"id":         "412904da-43ce-4aab-b79b-beea21c5d87a",
		"document":   "33333333",
		"nombre":     "deybin yoni gil perez",
		"address":    "av. general cordoba",
		"birthdate":  "04/04/1994",
		"age":        31,
		"amount":     "15000.00",
		"credits":    40,
		"passwords":  "Deybin04",
		"key_secret": "hola soy Deybin",
	}

	crud := pgorm.SqlExecMultiple{}
	crud.New(database)

	trSucursal := crud.SetInfo(new(tables.Schema_Models).New(), dataInsert...)
	trAlmacen := crud.SetInfo(new(tables.Schema_Models2).New(), dataInsert2)

	if err := trAlmacen.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := trSucursal.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := crud.Exec(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Multiple_transaction(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")
	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"l_sucu": "sucursal de prueba",
		"l_dire": "indefinido",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
		"l_dire": "indefinido",
	})
	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
		"l_dire": "indefinido",
	})

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	// schema, tableName := table.GetSucursal()
	// schemaAlma, tableNameAlma := table.GetAlmacen()
	crud := new(pgorm.SqlExecMultiple).New(database)
	TransactionAlmacen := crud.New(database).SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma...)
	TransactionSucursal := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert...)

	err := TransactionAlmacen.Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	err = TransactionSucursal.Insert()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.ExecTransaction(TransactionAlmacen)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.ExecTransaction(TransactionSucursal)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Commit()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
}

func TestCRUD_Multiple_Set_Transaction(t *testing.T) {

	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")

	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})

	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "004",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "003",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})
	dataInsertAlma = append(dataInsertAlma, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	dataInsert2 := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})

	dataInsertAlma2 := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	crud := pgorm.SqlExecMultiple{}
	crud.New(database)

	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	crud2 := pgorm.SqlExecMultiple{}
	crud2.New(database)

	if err := crud2.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert2...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	if err := crud2.SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma2...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if _, err := crud.SetTransaction(crud2.GetTransactions()...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := crud.Exec(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Multiple_Set_Single(t *testing.T) {

	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")

	dataInsert := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "003",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})

	dataInsert = append(dataInsert, map[string]interface{}{
		"c_sucu": "004",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})

	dataInsertAlma := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "003",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})
	dataInsertAlma = append(dataInsertAlma, map[string]interface{}{
		"c_sucu": "004",
		"c_alma": "001",
		"l_alma": "sucursal de prueba",
	})

	dataInsert2 := append([]map[string]interface{}{}, map[string]interface{}{
		"c_sucu": "005",
		"l_sucu": "sucursal de prueba",
		"l_dire": "sin información",
	})

	crud := pgorm.SqlExecMultiple{}
	crud.New(database)

	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsert...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	if err := crud.SetInfo(new(tables.Schema_Sucursal).New(), dataInsertAlma...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	crud2 := pgorm.SqlExecSingle{}
	crud2.New(new(tables.Schema_Sucursal).New(), dataInsert2...)

	if err := crud2.New(new(tables.Schema_Sucursal).New(), dataInsert2...).Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if _, err := crud.SetSqlSingle(crud2); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := crud.Exec(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_Single_ADD_SUM(t *testing.T) {
	godotenv.Load()
	database := os.Getenv("ENV_TES_DB")
	dataUpdate := map[string]interface{}{
		"credits": uint64(1),
		"where":   map[string]interface{}{"id": "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
	}

	crud := pgorm.SqlExecSingle{}
	err := crud.New(new(tables.Schema_Models).New(), dataUpdate).Update()
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	err = crud.Exec(database, true)
	if err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}
