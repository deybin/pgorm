package test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/deybin/pgorm"
	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/clause"

	tables "github.com/deybin/pgorm/test/table"
)

func Test_QueryFullString(t *testing.T) {
	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()
	data, err := pgorm.ExecQuery[map[string]any](db, context.Background(), pgorm.NewQuery().WorkQueryFull("SELECT * FROM models WHERE id=$1", "550e8400-e29b-41d4-a716-446655440001"))
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_Query__Sintaxis(t *testing.T) {

	var querySql = pgorm.NewQuery()

	queryString := querySql.Select().From(tables.Models{}.Name()).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("age", clause.I, 31).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age = $1" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("document", clause.I, "3345431").And("age", clause.I, 31).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE document = $1 AND age = $2" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("nombre", clause.I, "Juan").Or("address", clause.I, "Lima").String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE nombre = $1 OR address = $2" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("nombre", clause.I, "Juan").Or("address", clause.D, "Lima").And("age", clause.I, 31).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE nombre = $1 OR address <> $2 AND age = $3" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("atCreate", clause.BETWEEN, []any{"2025-01-01", "2025-01-31"}).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE atCreate BETWEEN $1 AND $2" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("age", clause.IN, []any{18, 21, 30}).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age IN ($1, $2, $3)" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("models.age", clause.I, 31).Join(clause.INNER, tables.Models2{}.Name(), "models.document=models2.document").String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models INNER JOIN models2 ON models.document=models2.document WHERE models.age = $1" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("age", clause.I, 31).Top(5).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age = $1 LIMIT 5" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.Models{}.Name()).Where("age", clause.I, 31).Limit(5, 10).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age = $1 LIMIT 5 OFFSET 10" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()
}

func Test_Query__Response(t *testing.T) {

	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})

	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()
	var querySql = pgorm.NewQuery()

	data, err := pgorm.ExecQuery[[]tables.Models](db, context.Background(), querySql.From(tables.Models{}.Name()).Select())
	if err != nil {
		t.Errorf("query inesperado: %q", err)
		return
	}
	fmt.Println(data)

	one, errOne := pgorm.ExecQuery[tables.Models](db, context.Background(), querySql.From(tables.Models{}.Name()).Select().Where("document", clause.I, "12345678903"))
	if errOne != nil {
		t.Errorf("query inesperado: %q", errOne)
		return
	}
	fmt.Println(one)

	type response struct {
		Document string `json:"document"`
		Nombre   string `json:"name"`
	}

	one2, errOne2 := pgorm.ExecQuery[response](db, context.Background(), querySql.From(tables.Models{}.Name()).Select("document, nombre").Where("document", clause.I, "12345678903"))
	if errOne2 != nil {
		t.Errorf("query inesperado: %q", errOne2)
		return
	}
	fmt.Println(one2)

	value, errValue := pgorm.ExecQuery[tables.Models](db, context.Background(), querySql.From(tables.Models{}.Name()).Select())
	if errValue != nil {
		t.Errorf("query inesperado: %q", errValue)
		return
	}
	fmt.Println(value)
}

func Test_Query__ResponseWithSchema(t *testing.T) {

	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})

	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()
	var querySql = pgorm.NewQuery()

	data, err := pgorm.ExecQueryWithSchema[[]tables.Models](db, "public", context.Background(), querySql.From(tables.Models{}.Name()).Select())
	if err != nil {
		t.Errorf("query inesperado: %q", err)
		return
	}
	fmt.Println(data)

	one, errOne := pgorm.ExecQueryWithSchema[tables.Models](db, "public", context.Background(), querySql.From(tables.Models{}.Name()).Select().Where("document", clause.I, "12345678903"))
	if errOne != nil {
		t.Errorf("query inesperado: %q", errOne)
		return
	}
	fmt.Println(one)

	type response struct {
		Document string `json:"document"`
		Nombre   string `json:"name"`
	}

	one2, errOne2 := pgorm.ExecQueryWithSchema[response](db, "public", context.Background(), querySql.From(tables.Models{}.Name()).Select("document, nombre").Where("document", clause.I, "12345678903"))
	if errOne2 != nil {
		t.Errorf("query inesperado: %q", errOne2)
		return
	}
	fmt.Println(one2)

	value, errValue := pgorm.ExecQueryWithSchema[tables.Models](db, "public", context.Background(), querySql.From(tables.Models{}.Name()).Select())
	if errValue != nil {
		t.Errorf("query inesperado: %q", errValue)
		return
	}
	fmt.Println(value)
}

func Test_Query__ResponseWithSchemaInContext(t *testing.T) {

	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	ctx := context.WithValue(context.Background(), pgorm.SchemaId, "public")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()
	var querySql = pgorm.NewQuery()

	data, err := pgorm.ExecQuery[[]tables.Models](db, ctx, querySql.From(tables.Models{}.Name()).Select())
	if err != nil {
		t.Errorf("query inesperado: %q", err)
		return
	}
	fmt.Println(data)

	one, errOne := pgorm.ExecQuery[tables.Models](db, ctx, querySql.From(tables.Models{}.Name()).Select().Where("document", clause.I, "12345678903"))
	if errOne != nil {
		t.Errorf("query inesperado: %q", errOne)
		return
	}
	fmt.Println(one)

	type response struct {
		Document string `json:"document"`
		Nombre   string `json:"name"`
	}

	one2, errOne2 := pgorm.ExecQuery[response](db, ctx, querySql.From(tables.Models{}.Name()).Select("document, nombre").Where("document", clause.I, "12345678903"))
	if errOne2 != nil {
		t.Errorf("query inesperado: %q", errOne2)
		return
	}
	fmt.Println(one2)

	value, errValue := pgorm.ExecQuery[tables.Models](db, ctx, querySql.From(tables.Models{}.Name()).Select())
	if errValue != nil {
		t.Errorf("query inesperado: %q", errValue)
		return
	}
	fmt.Println(value)
}

func Test_QueryProcedure__Error(t *testing.T) {
	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()

	id := "94f596c6-0dd1-4eae-871a-37d1dac81b28"
	err = pgorm.ExecProcedure(db, context.Background(), pgorm.NewQuery().WorkQueryFull(`
		DO $$
			DECLARE
						_num_credits int;
			BEGIN
				SELECT INTO _num_credits
						credits 
				FROM  models WHERE id=$1;

				_num_credits=_num_credits+1;
						
				UPDATE models SET credits=_num_credits	WHERE id=$2;
						
			END;
		$$ LANGUAGE plpgsql;
	`, id, id))

	if err != nil {
		if err.Error() != "mismatched param and argument count" {
			t.Errorf("no se esperaba este error: %s", err.Error())
			return
		}

	}

}

func Test_Query_Contexto(t *testing.T) {

	db, err := adapters.NewPool(adapters.ConfigPgxAdapter{})
	if err != nil {
		fmt.Println(err)
	}
	defer db.Pool().Close()
	var querySql = pgorm.NewQuery()
	ctx := context.Background()

	_, err = pgorm.ExecQuery[[]map[string]any](db, ctx, querySql.From(tables.Models{}.Name()).Select("document,id"))

	fmt.Println("session permanecerá por 10s primera vez!!!")
	time.Sleep(10 * time.Second)

	if err != nil {
		t.Errorf("No se esperaba este error: %v", err)
	}

	_, err = pgorm.ExecQuery[[]map[string]any](db, ctx, querySql.From(tables.Models{}.Name()).Select("nombre,address"))

	fmt.Println("session permanecerá por 10s mas")
	time.Sleep(10 * time.Second)

	querySql.Close()
	fmt.Println("session ya debe estar cerrado")
	time.Sleep(5 * time.Second)

	if err != nil {
		t.Errorf("No se esperaba este error: %v", err)
	}
}
