package migrator

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"strings"
	"unicode"

	"github.com/deybin/pgorm/internal/adapters"
)

// Table
// ParseInsert
// ParseUpdate
// ParseDelete

func GenerateSchemaFile(database string, table string) {

	dataTypeCollection := map[string]string{
		"uuid":      "string",
		"varchar":   "string",
		"char":      "string",
		"text":      "string",
		"int4":      "uint64",
		"float8":    "float64",
		"real":      "float64",
		"date":      "*time.Time",
		"timestamp": "*time.Time",
	}

	basePath := "./tmp/"

	resultTables := consultar("SELECT table_name FROM Information_Schema.TABLES WHERE table_name='"+table+"'", database)

	for _, v := range resultTables {

		tableName := toPascalCase(v["table_name"].(string))
		temp_table_name := strings.Split(tableName, "_")
		code_struct := "package tables\n"
		table := temp_table_name[0]
		nameSchema := "Schema" + table
		backticks := "`"
		code_struct += fmt.Sprintf(`
			type %s struct {
				table %s
			}
			`, nameSchema, tableName)

		query_sql := fmt.Sprintf("SELECT * FROM Information_Schema.Columns where  TABLE_NAME='%s'", v["table_name"].(string))
		resultColumns := consultar(query_sql, database)
		// fmt.Println(resultColumns)
		var codigo_schema string
		fmt.Println("tabla:", tableName, "   columnas:", len(resultColumns))
		var structTable []string
		for _, column := range resultColumns {
			var tagColumn string
			tagColumn += toPascalCase(column["column_name"].(string))
			tagColumn += fmt.Sprintf(` %s %s json:"%s" tag:"%s"  validate:"required"`, dataTypeCollection[column["udt_name"].(string)], backticks, column["column_name"], column["column_name"])

			var validateType []string
			if column["udt_name"] == "varchar" || column["udt_name"] == "char" || column["udt_name"] == "text" {
				if column["udt_name"] != "text" {
					max_length := int(column["character_maximum_length"].(int32))
					if max_length != 36 {
						if column["udt_name"] == "char" {
							validateType = append(validateType, fmt.Sprintf(`min=%d`, max_length))
						} else {
							validateType = append(validateType, fmt.Sprintf(`min=%d`, int(math.Round(float64(max_length)*0.1))))
						}
						validateType = append(validateType, fmt.Sprintf(`max=%d`, max_length))
					}
				}
				validateType = append(validateType, "case=lowercase")
			} else if column["data_type"] == "integer" {
				validateType = append(validateType, "max=10")
			} else if column["udt_name"] == "float8" || column["udt_name"] == "real" {
			} else if column["data_type"] == "date" {
			}

			tagColumn += fmt.Sprintf(` validateType:"%s" %s`, strings.Join(validateType, ","), backticks)
			structTable = append(structTable, tagColumn)

		}
		// fmt.Println(structTable)

		code_struct += fmt.Sprintf(`
type %s struct {
%s
}

func (s %s) Name() string {
	t := reflect.TypeOf(s)
	return strings.ToLower(t.Name())
}

func (s %s) Columns() []string {
	return migrator.EntityColumns(s)
}

func (s %s) Values() []any {
	return migrator.EntityValues(s)
}

func (s *%s) Table() migrator.Entity {
	return s.table
}

func (s *%s) ParseInsert() []migrator.Fields {
	return migrator.GenerateSchema(s.table, migrator.INSERT)
}

func (s *%s) ParseUpdate() []migrator.Fields {
	return migrator.GenerateSchema(s.table, migrator.UPDATE)
}

func (s *%s) ParseDelete() []migrator.Fields {
	return migrator.GenerateSchema(s.table, migrator.DELETE)
}
`, tableName, strings.Join(structTable, "\n"), tableName, tableName, tableName, nameSchema, nameSchema, nameSchema, nameSchema)

		codigo_schema += code_struct
		texto := []byte(code_struct)
		errs := os.WriteFile(fmt.Sprintf("%s%s.go", basePath, table), texto, 0644)
		if errs != nil {
			log.Fatal(errs)
		}

	}
}

func toPascalCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		if len(word) > 0 {
			// Capitaliza la primera letra de cada palabra
			words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
		}
	}
	return strings.Join(words, "_")
}

func consultar(query string, database string) []map[string]interface{} {
	db, err := adapters.NewPoolWithConfig(adapters.ConfigPgxAdapter{Database: database})
	ctx := context.Background()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// fmt.Println("Query: ", errPool)

	rows, err := db.Pool().Query(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	fieldDescription := rows.FieldDescriptions()
	cols := make([]string, len(fieldDescription))
	for i, fd := range fieldDescription {
		cols[i] = string(fd.Name)
	}

	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			log.Fatal(err)
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			l := *val

			if l != nil {
				if strings.Contains(reflect.TypeOf(l).String(), "[]uint8") {
					m[colName] = string(l.([]uint8))
				} else {
					m[colName] = l
				}
			} else {
				m[colName] = l
			}

		}

		result = append(result, m)
	}
	return result
}

/*
type tableModel struct {
	Document        string    `json:"document" field:"document" tag:"document" validate:"primaryKey,required" validateType:"case=lowercase,min=7,max=11,expr=number"`
	Nombre          string    `json:"nombre" field:"nombre" tag:"nombre" validate:"required, default=nombre" validateType:"case=lowercase,min=3,max=50,encrypt, cipher"`
	Password        string    `json:"password" field:"password" tag:"password" validate:"required" validateType:"encrypt"`
	Code_Secret     string    `json:"code_secret" field:"code_secret" tag:"Código secreto" validate:"required" validateType:"cipher"`
	FechaNacimiento time.Time `json:"fecha_nacimiento" filed:"fecha_nacimiento" tag:"fecha de nacimiento" validate:"required"`
	Age             uint64    `json:"age" field:"age" tag:"Edad" validate:"required" validateType:"max=80"`
	Amount          float64   `json:"amount" field:"amount" tag:"Monto dinerario" validate:"required,update,sum" validateType:"negativo,porcentaje,menor=40.00, mayor=50.00"`
	Credits         int64     `json:"credits" field:"credits" tag:"Créditos" validate:"required" validateType:"negativo,min=40, max=5"`
}
*/

func EntityColumns(s Entity) []string {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var fields []string
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !v.Field(i).CanInterface() {
			continue
		}
		name := field.Tag.Get("tag")
		if name == "" {
			name = field.Name
		}
		fields = append(fields, name)
	}
	return fields
}

func EntityValues(s Entity) []any {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	var values []any

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// ignorar campos no exportados
		if !v.Field(i).CanInterface() {
			continue
		}

		name := field.Tag.Get("tag")
		if name == "" {
			name = field.Name
		}

		values = append(values, v.Field(i).Interface())
	}

	return values
}
