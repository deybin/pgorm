package test

import (
	"testing"

	"github.com/deybin/pgorm/internal"
)

func TestConnectionDatabaseMaster(t *testing.T) {
	_, err := new(internal.Connection).New("").PoolMaster()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
}

func TestConnectionDatabaseVariable(t *testing.T) {
	_, err := new(internal.Connection).New("new_capital").NewPool()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}
}

func TestConnectionDatabaseVariable__Error(t *testing.T) {
	_, err := new(internal.Connection).New("new_capitala").NewPool()

	if err != nil {
		if err.Error() != "conexión error: base de datos no existe" {
			t.Errorf("este error no se esperaba: %s", err.Error())

		}
	} else {
		t.Errorf("se esperaba un error")
	}
}

func TestConnection__Close(t *testing.T) {
	cnn, err := new(internal.Connection).New("new_capital").NewPool()
	if err != nil {
		t.Errorf("no se esperaba este error: %s", err.Error())
		return
	}

	cnn.Close()

}
