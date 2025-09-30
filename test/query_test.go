package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deybin/pgorm"
)

func Test_QueryFullString(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetQueryString("SELECT * FROM models WHERE id=$1", "369d48ab-d881-4e22-9f14-8e9af67da9aa").Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__All(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__One(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Exec().One()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Text(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Exec().Text("id")
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Where(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Where("credits", pgorm.I, 2).Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__Top(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Top(2).Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(len(data))
}

func Test_QuerySelect__Limit(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Limit(3, 1).Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(len(data))
}

func Test_QuerySelect__WhereOr(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Where("credits", pgorm.I, 43).Or("id", pgorm.I, "94f596c6-0dd1-4eae-871a-37d1dac81b28").Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QuerySelect__WhereAnd(t *testing.T) {
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Where("credits", pgorm.I, 43).And("id", pgorm.I, "94f596c6-0dd1-4eae-871a-37d1dac81b28").Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println(data)
}

func Test_QueryProcedure__Error(t *testing.T) {
	id := "94f596c6-0dd1-4eae-871a-37d1dac81b28"
	err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetQueryString(`
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
	queryBuilder := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"})
	_, err := queryBuilder.SetTable("models").Select("document,id").ExecCtx().All()

	fmt.Println("session permanecerá por 10s primera vez!!!")
	time.Sleep(10 * time.Second)

	if err != nil {
		t.Errorf("No se esperaba este error: %v", err)
	}

	_, err = queryBuilder.SetTable("models").Select("nombre,address").ExecCtx().All()

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
	data, err := new(pgorm.Query).New(pgorm.QConfig{Database: "new_capital"}).SetTable("models").Select().Exec().All()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
	fmt.Println("10s para ver si existe conexión activa!!!")
	time.Sleep(10 * time.Second)
	fmt.Println(data)
}
