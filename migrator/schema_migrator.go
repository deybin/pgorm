package migrator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Entity interface {
	Name() string
	Columns() []string
	Values() []any
}

// Schema es una interfaz que define métodos para obtener información sobre el esquema de una tabla en una base de datos.
type Schema interface {
	Table() Entity
	ParseInsert() []Fields
	ParseUpdate() []Fields
	ParseDelete() []Fields
}

type (
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
	NameOriginal         string   //Nombre del campo original
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

type Where struct {
	Clause    string
	Condition string
	Field     string
	Value     any
}

type ModelsFilter[T any] struct {
	Data       T
	Conditions []Where
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

type Actions uint

const (
	NONE Actions = iota
	INSERT
	UPDATE
	DELETE
)

type EntityUpdate struct {
	Entity
	Conditions []Where
}

func GenerateSchema(data any, action Actions) []Fields {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	var schema []Fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		var structSchema Fields
		structSchema.Name = strings.ToLower(field.Name)
		structSchema.NameOriginal = field.Name
		structSchema.Description = field.Tag.Get("tag")

		if structSchema.Name == "atcreate" {
			structSchema.Default = time.Now().UTC()
		}

		validateTag := field.Tag.Get("validate")
		structSchema.PrimaryKey = strings.Contains(validateTag, "primaryKey")
		structSchema.Where = strings.Contains(validateTag, "where")

		if action == DELETE {
			if !structSchema.PrimaryKey && !structSchema.Where {
				continue
			}
		}

		structSchema.Update = strings.Contains(validateTag, "update")
		if action == UPDATE {
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
		rules := strings.Split(validateTypeTag, ";")
		switch structSchema.Type {
		case String:
			schemaType := TypeStrings{}
			for _, rule := range rules {
				if after, ok := strings.CutPrefix(rule, "min="); ok {
					min, _ := strconv.Atoi(after)
					schemaType.Min = min
				}
				if after, ok := strings.CutPrefix(rule, "max="); ok {
					max, _ := strconv.Atoi(after)
					schemaType.Max = max
				}
				if after, ok := strings.CutPrefix(rule, "case="); ok {
					switch after {
					case "lowercase":
						schemaType.LowerCase = true
					case "uppercase":
						schemaType.UpperCase = true
					}
				}
				if strings.Contains(rule, "encrypt") {
					schemaType.Encriptar = true
				}
				if strings.Contains(rule, "cipher") {
					schemaType.Cifrar = true
				}
				if after, ok := strings.CutPrefix(rule, "expr="); ok {
					schemaType.Expr = regexp.MustCompile(after)

				}
			}
			structSchema.ValidateType = schemaType
		case Float:
			schemaType := TypeFloat64{}
			for _, rule := range rules {
				if after, ok := strings.CutPrefix(rule, "menor="); ok {
					menor, _ := strconv.ParseFloat(after, 64)
					schemaType.Menor = menor
				} else if after0, ok0 := strings.CutPrefix(rule, "mayor="); ok0 {
					mayor, _ := strconv.ParseFloat(after0, 64)
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
				if after, ok := strings.CutPrefix(rule, "min="); ok {
					min, _ := strconv.Atoi(after)
					schemaType.Min = int64(min)
				} else if after0, ok0 := strings.CutPrefix(rule, "max="); ok0 {
					max, _ := strconv.Atoi(after0)
					schemaType.Max = int64(max)
				} else if strings.Contains(rule, "negativo") {
					schemaType.Negativo = true
				}
			}
			structSchema.ValidateType = schemaType
		case Uint:
			schemaType := TypeUint64{}
			for _, rule := range rules {
				if after, ok := strings.CutPrefix(rule, "max="); ok {
					max, _ := strconv.Atoi(after)
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
