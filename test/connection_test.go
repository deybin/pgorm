package test

import (
	"testing"

	"github.com/deybin/pgorm"
)

func TestConnectionDatabaseMaster(t *testing.T) {
	_, err := new(pgorm.Connection).New("").PoolMaster()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
}

func TestConnectionDatabaseVariable(t *testing.T) {
	_, err := new(pgorm.Connection).New("new_capital").Pool()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
}

func TestConnectionDatabaseVariable__Error(t *testing.T) {
	_, err := new(pgorm.Connection).New("new_capitala").Pool()

	if err != nil {
		if err.Error() != "conexión error: base de datos no existe" {
			t.Errorf("este error no se esperaba: %s", err.Error())

		}
	} else {
		t.Errorf("se esperaba un error")
	}
}

func TestConnection__Close(t *testing.T) {
	cnn, err := new(pgorm.Connection).New("new_capital").Pool()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}

	cnn.Close()

}
