package migrator

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/deybin/pgorm/internal/configs"
	"github.com/deybin/pgorm/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func CheckInsertGeneric(schema []Fields, tabla_map Entity) (map[string]any, error) {
	var errs []string
	data := make(map[string]any)
	for _, item := range schema {
		v := reflect.ValueOf(tabla_map)
		// fmt.Println("ERRName:", item.NameOriginal)
		val := v.FieldByName(item.NameOriginal)
		// fmt.Println("ERR:", val)
		isNil := val.IsZero()
		// fmt.Println(isNil, item.NameOriginal)
		defaultIsNil := item.Default == nil
		if !isNil {
			if _, ok := val.Interface().(string); ok {
				if strings.TrimSpace(val.String()) == "" {
					isNil = true
				}
			}
		}

		if !isNil {

			val, err := validaciones(item, val.Interface())
			if err == nil {
				data[item.Name] = val
			} else {
				errs = append(errs, fmt.Sprintf("Se encontró fallas al validar el campo %s: (%s)", item.Description, err.Error()))
			}
		} else {
			if !defaultIsNil {
				data[item.Name] = item.Default
			} else {
				if item.Required || item.PrimaryKey {
					errs = append(errs, fmt.Sprintf("El campo %s: (es Requerido)", item.Description))
				}
			}
		}

	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, ", "))
	} else {
		return data, nil
	}
}

func CheckUpdateGeneric(schemas []Fields, tabla_map any) (map[string]any, error) {
	var errs []string
	data := make(map[string]any)
	for _, item := range schemas {
		if item.NameOriginal == "Conditions" {
			continue
		}
		v := reflect.ValueOf(tabla_map)
		val := v.FieldByName(item.NameOriginal)
		isNil := val.IsValid()
		if isNil {
			isNil = val.IsZero()
		}

		// fmt.Println(isNil, ":", val.Interface(), ":", item.NameOriginal)
		if !isNil {
			if item.Update {
				value := val.Interface()

				if x, ok := value.(string); ok {
					if strings.TrimSpace(x) == "" {
						if !item.Empty {
							errs = append(errs, fmt.Sprintf("El campo %s no puede estar vació\n", item.Description))
						}
					}
				}

				var val interface{}
				val, err := validaciones(item, value)

				if err == nil {
					keyName := item.Name

					switch item.ArithmeticOperations {
					case Sum:
						keyName = fmt.Sprintf("ADD_%s_SUMA", item.Name)
					case Subtraction:
						keyName = fmt.Sprintf("ADD_%s_SUBTRACTION", item.Name)
					case Multiply:
						keyName = fmt.Sprintf("ADD_%s_MULTIPLY", item.Name)
					case Divide:
						keyName = fmt.Sprintf("ADD_%s_DIVIDE", item.Name)
					}

					data[keyName] = val
				} else {
					errs = append(errs, fmt.Sprintf("Se encontró fallas al validar el campo %s \n %s\n", item.Description, err.Error()))
				}
			} else {
				errs = append(errs, fmt.Sprintf("El campo %s no puede ser modificado\n", item.Description))
			}
		}
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, ", "))
	} else {
		return data, nil
	}
}

func CheckWhereGeneric(schemas []Fields, table_where ...Where) ([]Where, error) {
	var errs []string
	usePrimaryKey := false
	fieldNotExistLent := 0
	for _, item := range table_where {
		fieldExist := false
		for _, v := range schemas {
			// fmt.Println(item.Field, " : ", v.Name)
			if item.Field == v.Name {
				if v.Where || v.PrimaryKey {
					if v.PrimaryKey {
						usePrimaryKey = true
					}
					if x, ok := item.Value.(string); ok {
						if strings.TrimSpace(x) == "" {
							errs = append(errs, fmt.Sprintf("El campo %s no puede ser utilizado de esta forma\n", v.Description))
						}

					}

				} else {
					errs = append(errs, fmt.Sprintf("El campo %s no puede ser utilizado de esta forma\n", v.Description))
				}
				fieldExist = true
				continue
			}

		}
		if !fieldExist {
			fieldNotExistLent++
		}

	}

	if fieldNotExistLent > 0 {
		errs = append(errs, "Uno o más campos enviados no son válidos.")
	}

	existePrimaryKey := slices.ContainsFunc(schemas, func(u Fields) bool {
		return u.PrimaryKey == true
	})

	if !usePrimaryKey || !existePrimaryKey {
		errs = append(errs, "Existen campos obligatorios sin information")
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, ", "))
	} else {
		return table_where, nil
	}
}

func validaciones(item Fields, value any) (any, error) {
	var val any
	var err error

	switch v := value.(type) {
	case string:
		val, err = caseString(v, item.ValidateType.(TypeStrings))
	case *string:
		val, err = caseString(*v, item.ValidateType.(TypeStrings))
	case float64:
		val, err = caseFloat(v, item.ValidateType.(TypeFloat64))
	case *float64:
		val, err = caseFloat(*v, item.ValidateType.(TypeFloat64))
	case uint64:
		val, err = caseUint(v, item.ValidateType.(TypeUint64))
	case *uint64:
		val, err = caseUint(*v, item.ValidateType.(TypeUint64))
	case int64:
		val, err = caseInt(v, item.ValidateType.(TypeInt64))
	case *int64:
		val, err = caseInt(*v, item.ValidateType.(TypeInt64))
	case *time.Time:
		val = v // Los tiempos se suelen manejar como punteros directamente
	case bool:
		val = v
	case *bool:
		val = *v // Ahora es seguro obtener el contenido

	default:
		val, err = nil, errors.New("tipo de dato no asignado")
	}

	return val, err
}

func caseString(value string, schema TypeStrings) (string, error) {
	value = strings.TrimSpace(value)
	if schema.Expr != nil {
		if !schema.Expr.MatchString(value) {
			return "", errors.New("no cumple con las características")
		}
	}

	if schema.Encriptar {
		result, err := bcrypt.GenerateFromPassword([]byte(value), 13)
		if err != nil {
			return value, err
		}
		value = string(result)
		return value, nil
	}

	if schema.Cifrar {
		hash, err := utils.AesEncrypt_PHP([]byte(value), configs.KeyCrypto())
		if err != nil {
			return value, err
		}
		value = hash
		return value, nil
	}

	if schema.Min > 0 {
		if len(value) < schema.Min {
			return "", fmt.Errorf("no Cumple los caracteres mínimos que debe tener (%v)", schema.Min)
		}
	}

	if schema.Max > 0 {
		if len(value) > schema.Max {
			return "", fmt.Errorf("no Cumple los caracteres máximos que debe tener (%v)", schema.Max)
		}
	}

	if schema.UpperCase {
		value = strings.ToUpper(value)
	} else if schema.LowerCase {
		value = strings.ToLower(value)
	}

	return value, nil
}

func caseDate(value *time.Time, _ TypeDate) (*time.Time, error) {
	return value, nil
}

func caseFloat(value float64, schema TypeFloat64) (float64, error) {
	var err []string
	if schema.Menor != 0 {
		if value <= schema.Menor {
			err = append(err, fmt.Sprintf("No puede se menor a %f", schema.Menor))
		}
	}
	if schema.Mayor != 0 {
		if value >= schema.Mayor {
			err = append(err, fmt.Sprintf("No puede se mayor a %f", schema.Mayor))
		}
	}
	if !schema.Negativo {
		if value < float64(0) {
			err = append(err, "No puede ser negativo")
		}
	}
	if schema.Porcentaje {
		value = value / float64(100)
	}
	if len(err) > 0 {
		return 0, errors.New(strings.Join(err, ", "))
	} else {
		return value, nil
	}
}

func caseInt(value int64, schema TypeInt64) (int64, error) {
	var err []string
	if !schema.Negativo {
		if value < int64(0) {
			err = append(err, "No puede ser negativo")
		}
	}
	if schema.Min != 0 {
		if value < schema.Min {
			err = append(err, fmt.Sprintf("No puede se menor a %d", schema.Min))
		}
	}
	if schema.Max != 0 {
		if value > schema.Max {
			err = append(err, fmt.Sprintf("No puede se mayor a %d", schema.Max))
		}
	}
	if len(err) > 0 {
		return int64(0), errors.New(strings.Join(err, ", "))
	} else {
		return value, nil
	}
}

func caseUint(value uint64, t TypeUint64) (uint64, error) {
	if t.Max > 0 {
		if value > t.Max {
			return 0, errors.New("no esta en el rango permitido")
		}
	}
	return value, nil
}
