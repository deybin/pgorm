package pgorm

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Schema es una interfaz que define métodos para obtener información sobre el esquema de una tabla en una base de datos.
type Schema interface {
	GetTableName() string
	GetSchemaInsert() []Fields
	GetSchemaUpdate() []Fields
	GetSchemaDelete() []Fields
}

type (
	// DataType basicGORM data type
	DataType string
)

const (
	Bool   DataType = "bool"
	Int    DataType = "int64"
	Uint   DataType = "uint64"
	Float  DataType = "float64"
	String DataType = "string"
	Time   DataType = "time.Time"
	Bytes  DataType = "bytes"
)

type ArithmeticOperations uint

const (
	None ArithmeticOperations = iota
	Sum
	Subtraction
	Multiply
	Divide
)

/*
Las etiquetas de structure o también llamado etiquetas de campo estos metadatos serán los siguientes según el tipo de dato
*/
type Fields struct {
	Name                 string   //Nombre del campo
	Description          string   //Descripción del campo
	Type                 DataType //A bajo nivel es un string donde se especifica de que tipo sera el campo
	ArithmeticOperations ArithmeticOperations
	Required             bool        //Si el valor para inserción de este campo es requerido o no
	PrimaryKey           bool        //Si el campo es primary key entonces es obligatorio este campo para insert,update y delete
	Where                bool        //El campo puede ser utilizado para filtrar al utilizar el update y delete
	Update               bool        //El campo puede ser modificado
	Default              interface{} //Valor por defecto que se tomara si no se le valor al campo, el tipo del valor debe de ser igual al Type del campo
	Empty                bool        //El campo aceptara valor vació si se realiza la actualización
	ValidateType         interface{} //Los datos serán validados mas a fondo mediante esta opción para eso se le debe de asignar los siguientes typo de struct: TypeStrings, TypeFloat64, TypeUint64 yTypeInt64
}

type TypeStrings struct {
	LowerCase bool           //Convierte en minúscula el valor del campo
	UpperCase bool           //Convierte en mayúscula el valor del campo
	Encriptar bool           //Crea un hash del valor del campo
	Cifrar    bool           //Cifra el valor del campo
	Date      bool           //Verifica que el valor del campo sea una fecha valida con formato dd/mm/yyyy
	Min       int            //Cuantos caracteres como mínimo debe de tener el valor del campo
	Max       int            //Cuantos caracteres como máximo debe de tener el valor del campo
	Expr      *regexp.Regexp //Expresión regular que debe cumplir el valor que almacenara el campo
}

type TypeDate struct {
}

type TypeFloat64 struct {
	Porcentaje bool    //Convierte el valor del campo en porcentaje
	Negativo   bool    //El campo aceptara valores negativos
	Menor      float64 //Los valores que aceptaran tienen que ser menores o igual que este metadato
	Mayor      float64 //Los valores que aceptaran tienen que ser mayores o igual  que este metadato
}

type TypeUint64 struct {
	Max uint64 //Hasta que valor aceptara que se almacene en este campo
}

type TypeInt64 struct {
	Max      int64 // Hasta que valor aceptara que se almacene en este campo
	Min      int64 // Valor como mínimo  aceptara que se almacene en este campo
	Negativo bool  // Rl campo aceptara valores negativos
}

type Regex interface {
	Letras(start int8, end int16) *regexp.Regexp
	Float() *regexp.Regexp
}

func Null() *regexp.Regexp {
	return regexp.MustCompile(``)
}

func Number() *regexp.Regexp {
	return regexp.MustCompile(`[0-9]{0,}$`)
}

type ActionType uint

const (
	ActionInsert ActionType = iota
	ActionUpdate
	ActionDelete
)

type SchemaExe struct{}

func (s SchemaExe) GenerateSchema(data any, action ActionType) []Fields {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	var schema []Fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		var structSchema Fields
		structSchema.Name = strings.ToLower(field.Name)
		structSchema.Description = field.Tag.Get("tag")

		if structSchema.Name == "atcreate" {
			structSchema.Default = Date{}.GetDateLocation()
		}

		validateTag := field.Tag.Get("validate")
		structSchema.PrimaryKey = strings.Contains(validateTag, "primaryKey")
		structSchema.Where = strings.Contains(validateTag, "where")

		if action == ActionDelete {
			if !structSchema.PrimaryKey && !structSchema.Where {
				continue
			}
		}

		structSchema.Update = strings.Contains(validateTag, "update")
		if action == ActionUpdate {
			if !structSchema.Update && !structSchema.PrimaryKey && !structSchema.Where {
				continue
			}

			if strings.Contains(validateTag, "sum") {
				structSchema.ArithmeticOperations = Sum
			} else if strings.Contains(validateTag, "subtraction") {
				structSchema.ArithmeticOperations = Subtraction
			} else if strings.Contains(validateTag, "multiply") {
				structSchema.ArithmeticOperations = Multiply
			} else if strings.Contains(validateTag, "divide") {
				structSchema.ArithmeticOperations = Divide
			}
		}

		structSchema.Required = strings.Contains(validateTag, "required")

		structSchema.Type = DataType(field.Type.String())
		//fmt.Println("DataType(field.Type.Name()):", DataType(field.Type.String()))
		if strings.Contains(validateTag, "default") {
			structSchema.Default = value.Interface()
		}
		validateTypeTag := field.Tag.Get("validateType")
		rules := strings.Split(validateTypeTag, ",")
		switch structSchema.Type {
		case String:
			schemaType := TypeStrings{}
			for _, rule := range rules {
				if strings.HasPrefix(rule, "min=") {
					min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
					schemaType.Min = min
				} else if strings.HasPrefix(rule, "max=") {
					max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
					schemaType.Max = max
				} else if strings.HasPrefix(rule, "case=") {
					typeCase := strings.TrimPrefix(rule, "case=")
					switch typeCase {
					case "lowercase":
						schemaType.LowerCase = true
					case "uppercase":
						schemaType.UpperCase = true
					}
				} else if strings.Contains(rule, "encrypt") {
					schemaType.Encriptar = true
				} else if strings.Contains(rule, "cipher") {
					schemaType.Cifrar = true
				} else if strings.HasPrefix(rule, "expr=") {
					typeCase := strings.TrimPrefix(rule, "expr=")
					if typeCase == "number" {
						schemaType.Expr = Number()
					}
				}
			}
			structSchema.ValidateType = schemaType
		case Float:
			schemaType := TypeFloat64{}
			for _, rule := range rules {
				if strings.HasPrefix(rule, "menor=") {
					menor, _ := strconv.ParseFloat(strings.TrimPrefix(rule, "menor="), 64)
					schemaType.Menor = menor
				} else if strings.HasPrefix(rule, "mayor=") {
					mayor, _ := strconv.ParseFloat(strings.TrimPrefix(rule, "mayor="), 64)
					schemaType.Mayor = mayor
				} else if strings.Contains(rule, "negative") {
					schemaType.Negativo = true
				} else if strings.Contains(rule, "porcentaje") {
					schemaType.Porcentaje = true
				}
			}
			structSchema.ValidateType = schemaType
		case Int:
			schemaType := TypeInt64{}
			for _, rule := range rules {
				if strings.HasPrefix(rule, "min=") {
					min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
					schemaType.Min = int64(min)
				} else if strings.HasPrefix(rule, "max=") {
					max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
					schemaType.Max = int64(max)
				} else if strings.Contains(rule, "negativo") {
					schemaType.Negativo = true
				}
			}
			structSchema.ValidateType = schemaType
		case Uint:
			schemaType := TypeUint64{}
			for _, rule := range rules {
				if strings.HasPrefix(rule, "max=") {
					max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
					schemaType.Max = uint64(max)
				}
			}
			structSchema.ValidateType = schemaType
		case Time:
			schemaType := TypeDate{}
			structSchema.ValidateType = schemaType
		}

		schema = append(schema, structSchema)
	}

	return schema
}

func (s SchemaExe) GenerateSchemaFile(database string, table string) {

	dataTypeCollection := map[string]string{
		"varchar":   "string",
		"char":      "string",
		"text":      "string",
		"int4":      "uint64",
		"float8":    "float64",
		"real":      "float64",
		"date":      "time.Time",
		"timestamp": "time.Time",
	}

	basePath := "./tmp/"

	resultTables := consultar("SELECT table_name FROM Information_Schema.TABLES WHERE table_name='"+table+"'", database)

	for _, v := range resultTables {

		tableName := toPascalCase(v["table_name"].(string))
		temp_table_name := strings.Split(tableName, "_")
		code_struct := "package tables\n"
		table := temp_table_name[0]
		nameSchema := "Schema_" + table
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

func (s *%s) New() *%s {
	return s
}

func (s *%s) GetTableName() string {
	t := reflect.TypeOf(s.table)
	return strings.ToLower(t.Name())
}

func (s *%s) GetSchemaInsert() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionInsert)
}

func (s *%s) GetSchemaUpdate() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionUpdate)
}

func (s *%s) GetSchemaDelete() []pgorm.Fields {
	return pgorm.SchemaExe{}.GenerateSchema(s.table, pgorm.ActionDelete)
}
`, tableName, strings.Join(structTable, "\n"), nameSchema, nameSchema, nameSchema, nameSchema, nameSchema, nameSchema)

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
	db, err := new(Connection).New(database).Pool()
	ctx := context.Background()

	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Query: ", errPool)

	rows, err := db.pool.Query(ctx, query)
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
