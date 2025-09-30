package test

import (
	"testing"

	"github.com/deybin/pgorm"
)

func TestModelsGenerateFile(t *testing.T) {

	pgorm.SchemaExe{}.GenerateSchemaFile("new_capital", "models")

	r := ""

	result := []map[string]interface{}{}
	if len(result) != 0 {
		t.Errorf("Se esperaba: %v, pero se obtuvo %v", result, r)
	}
}
