package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/deybin/pgorm"
	"github.com/deybin/pgorm/clause"
	tables "github.com/deybin/pgorm/test/table"
)

func Test_QueryFullString(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "test"}).WorkQueryFull("SELECT * FROM models WHERE id=$1", "369d48ab-d881-4e22-9f14-8e9af67da9aa").Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Sintaxis(t *testing.T) {
	var querySql = pgorm.NewQuery(pgorm.Config{Database: "test"})

	queryString := querySql.Select().From(tables.ModelsSchema{}).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("age", clause.I, 31).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age = $1" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("document", clause.I, "3345431").And("age", clause.I, 31).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE document = $1 AND age = $2" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("nombre", clause.I, "Juan").Or("address", clause.I, "Lima").String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE nombre = $1 OR address = $2" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("nombre", clause.I, "Juan").Or("address", clause.D, "Lima").And("age", clause.I, 31).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE nombre = $1 OR address <> $2 AND age = $3" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("atCreate", clause.BETWEEN, []interface{}{"2025-01-01", "2025-01-31"}).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE atCreate BETWEEN $1 AND $2" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("age", clause.IN, []interface{}{18, 21, 30}).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age IN ($1, $2, $3)" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("models.age", clause.I, 31).Join(clause.INNER, tables.Models2Schema{}, "models.document=models2.document").String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models INNER JOIN models2 ON models.document=models2.document WHERE models.age = $1" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("age", clause.I, 31).Top(5).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age = $1 LIMIT 5" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()

	queryString = querySql.Select().From(tables.ModelsSchema{}).Where("age", clause.I, 31).Limit(5, 10).String()
	if strings.TrimSpace(queryString) != "SELECT * FROM models WHERE age = $1 LIMIT 5 OFFSET 10" {
		t.Errorf("query inesperado: %q", queryString)
		return
	}
	fmt.Println("sintaxis OK: ", queryString)
	querySql.Reset()
}

func Test_QuerySelect__AllStruct(t *testing.T) {

	data, err := pgorm.G[tables.Models](pgorm.Config{Database: "test"}).From(tables.ModelsSchema{}).Select().Exec().All()
	if err != nil {
		t.Errorf("query inesperado: %q", err)
		return
	}
	fmt.Println(data)

	one, errOne := pgorm.G[tables.Models](pgorm.Config{Database: "test"}).From(tables.ModelsSchema{}).Select().Where("document", clause.I, "12345678903").Exec().One()
	if errOne != nil {
		t.Errorf("query inesperado: %q", errOne)
		return
	}
	fmt.Println(one)

	type response struct {
		Document string `json:"document"`
		Nombre   string `json:"name"`
	}

	one2, errOne2 := pgorm.G[response](pgorm.Config{Database: "test"}).From(tables.ModelsSchema{}).Select("document, nombre").Where("document", clause.I, "12345678903").Exec().One()
	if errOne2 != nil {
		t.Errorf("query inesperado: %q", errOne2)
		return
	}
	fmt.Println(one2)

	value, errValue := pgorm.G[tables.Models](pgorm.Config{Database: "test"}).From(tables.ModelsSchema{}).Select().Exec().Value(func(m tables.Models) any { return m.Nombre })
	if errValue != nil {
		t.Errorf("query inesperado: %q", errValue)
		return
	}
	fmt.Println(value)
}
func Test_QuerySelect__All(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "test"}).From(tables.ModelsSchema{}).Select().Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__One(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Exec().One()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Text(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Exec().Text("id")
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Where(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Where("credits", clause.I, 2).Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Top(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Top(2).Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(len(data))
}

func Test_QuerySelect__Limit(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Limit(3, 1).Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(len(data))
}

func Test_QuerySelect__WhereOr(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Where("credits", clause.I, 43).Or("id", clause.I, "94f596c6-0dd1-4eae-871a-37d1dac81b28").Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__WhereAnd(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Where("credits", clause.I, 43).And("id", clause.I, "94f596c6-0dd1-4eae-871a-37d1dac81b28").Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QueryProcedure__Error(t *testing.T) {
	id := "94f596c6-0dd1-4eae-871a-37d1dac81b28"
	err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).WorkQueryFull(`
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
	`, id, id).Procedure()
	if err != nil {
		if err.Error() != "mismatched param and argument count" {
			t.Errorf("no se esperaba este error: %s", err.Error())
			return
		}

	}

}

func Test_Query_Contexto(t *testing.T) {
	queryBuilder := pgorm.NewQuery(pgorm.Config{Database: "new_capital"})
	_, err := queryBuilder.From(tables.ModelsSchema{}).Select("document,id").ExecCtx().All()

	fmt.Println("session permanecerá por 10s primera vez!!!")
	time.Sleep(10 * time.Second)

	if err != nil {
		t.Errorf("No se esperaba este error: %v", err)
	}

	_, err = queryBuilder.From(tables.ModelsSchema{}).Select("nombre,address").ExecCtx().All()

	fmt.Println("session permanecerá por 10s mas")
	time.Sleep(10 * time.Second)

	queryBuilder.Close()
	fmt.Println("session ya debe estar cerrado")
	time.Sleep(5 * time.Second)

	if err != nil {
		t.Errorf("No se esperaba este error: %v", err)
	}
}

// prueba si la conexión es cerrada  ni bien obtenemos los datos
func Test_QuerySelect__SinContext(t *testing.T) {
	data, err := pgorm.NewQuery(pgorm.Config{Database: "new_capital"}).From(tables.ModelsSchema{}).Select().Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println("10s para ver si existe conexión activa!!!")
	time.Sleep(10 * time.Second)
	fmt.Println(data)
}
