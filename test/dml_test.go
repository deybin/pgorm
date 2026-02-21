package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deybin/pgorm"
	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/migrator"
	tables "github.com/deybin/pgorm/test/table"
)

// inserta,actualiza y elimina datos solo de una tabla
func TestCRUD_IUD_Single(t *testing.T) {
	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()

	birth := time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local)
	dataInsert := tables.Models{
		Id:         "369d48ab-d881-4e22-9f14-8e9af67da9aa",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        31,
		Amount:     17000.00,
		Credits:    40,
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_secret: "hola soy Deybin",
	}

	crudInsert := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{}, dataInsert)

	if err := crudInsert.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err = pgorm.TransactionExec(crudInsert); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	time.Sleep(10 * time.Second)

	dataUpdate := migrator.EntityUpdate{
		Entity: tables.Models{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []migrator.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	crudUpdate := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{}, dataUpdate)

	if err := crudUpdate.Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := pgorm.TransactionExec(crudUpdate); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	time.Sleep(10 * time.Second)

	dataDelete := []migrator.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	crudDelete := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{})

	if err := crudDelete.Delete(dataDelete...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	if err := pgorm.TransactionExec(crudDelete); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
}

func TestCRUD_IUD_MultipleSuccess(t *testing.T) {
	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()

	birth := time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local)
	dataInsert := tables.Models{
		Id:         "550e8400-e29b-41d4-a716-446655440011",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        31,
		Amount:     15000.00,
		Credits:    40,
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_secret: "hola soy Deybin",
	}

	dataInsert2 := tables.Models2{
		Id:         "550e8400-e29b-41d4-a716-446655440012",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        31,
		Amount:     15000.00,
		Credits:    40,
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_secret: "hola soy Deybin",
	}

	// println(dataInsert.GetSchemaInsert())
	crud := pgorm.NewSqlExecMultiples(db)

	table := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{}, dataInsert))

	if err := table.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	table2 := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema2{}, dataInsert2))
	if err := table2.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataUpdate2 := migrator.EntityUpdate{
		Entity: tables.Models2{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []migrator.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "550e8400-e29b-41d4-a716-446655440012"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	table3 := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema2{}, dataUpdate2))
	if err := table3.Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataUpdate := struct {
		migrator.Entity
		Conditions []migrator.Where
	}{
		Entity: tables.Models{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []migrator.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "550e8400-e29b-41d4-a716-446655440011"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	crudUpdate := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{}, dataUpdate)
	if err := crudUpdate.Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	crud.SetTransactions(crudUpdate)

	dataDelete1 := []migrator.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	table4 := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema2{}))
	if err := table4.Delete(dataDelete1...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataDelete := []migrator.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	crudDelete := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{})

	if err := crudDelete.Delete(dataDelete...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	crud.SetTransactions(crudDelete)

	if err := pgorm.TransactionMultiExec(crud); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

}

func TestCRUD_IUD_MultipleError(t *testing.T) {
	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()

	birth := time.Date(1994, 4, 4, 0, 0, 0, 0, time.Local)
	dataInsert := tables.Models{
		Id:         "550e8400-e29b-41d4-a716-446655440011",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        31,
		Amount:     15000.00,
		Credits:    40,
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_secret: "hola soy Deybin",
	}

	dataInsert2 := tables.Models2{
		Id:         "550e8400-e29b-41d4-a716-446655440012",
		Document:   "719780401",
		Nombre:     "deybin yoni gil perez",
		Address:    "av. general cordoba 427",
		Birthdate:  &birth,
		Age:        31,
		Amount:     15000.00,
		Credits:    40,
		Passwords:  "Deybin041",
		Email:      "Deybin.04@gmail.com",
		Key_secret: "hola soy Deybin",
	}

	crud := pgorm.NewSqlExecMultiples(db)

	table := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{}, dataInsert))

	if err := table.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	table2 := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema2{}, dataInsert2))
	if err := table2.Insert(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataUpdate2 := migrator.EntityUpdate{
		Entity: tables.Models2{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre", Document: "123453212123"},
		Conditions: []migrator.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "550e8400-e29b-41d4-a716-446655440012"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	table3 := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema2{}, dataUpdate2))
	if err := table3.Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataUpdate := struct {
		migrator.Entity
		Conditions []migrator.Where
	}{
		Entity: tables.Models{Email: "deybin_04@hotmail.com", Nombre: "Nuevo Nombre"},
		Conditions: []migrator.Where{
			{Clause: "WHERE", Condition: "=", Field: "id", Value: "550e8400-e29b-41d4-a716-446655440011"},
			{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
		},
	}

	crudUpdate := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{}, dataUpdate)
	if err := crudUpdate.Update(); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}
	crud.SetTransactions(crudUpdate)

	dataDelete1 := []migrator.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	table4 := crud.Set(pgorm.NewSqlExecSingles(db, &tables.ModelsSchema2{}))
	if err := table4.Delete(dataDelete1...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	dataDelete := []migrator.Where{
		{Clause: "WHERE", Condition: "=", Field: "id", Value: "369d48ab-d881-4e22-9f14-8e9af67da9aa"},
		{Clause: "OR", Condition: "=", Field: "document", Value: "719780401"},
	}
	crudDelete := pgorm.NewSqlExecSingles(db, &tables.ModelsSchema{})

	if err := crudDelete.Delete(dataDelete...); err != nil {
		t.Errorf("se esperaba este error: %s", err.Error())
		return
	}

	crud.SetTransactions(crudDelete)
	if err := pgorm.TransactionMultiExec(crud); err != nil {
		if err.Error() == "ERROR: value too long for type character varying(11) (SQLSTATE 22001)" {
			fmt.Println("Error OK!!!: ", err.Error())
			return
		}

	}
	t.Errorf("se esperaba un error")

}
