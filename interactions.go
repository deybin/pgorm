package pgorm

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func Query_Cross_Update(query string) string {

	query_regEx := regexp.MustCompile(`ADD_(.*?)_SUMA=`)
	array := query_regEx.FindAllStringSubmatch(query, -1)
	for _, v := range array {

		new := fmt.Sprintf("%s=%s+", v[1], v[1])

		query = strings.Replace(query, v[0], new, -1)
	}

	return query
}

// recibe un valor interface que no se reconoce su tipo y devuelve un string
func InterfaceToString(params ...interface{}) string {
	typeValue := reflect.TypeOf(params[0]).String()
	value := params[0]
	valueReturn := ""
	if strings.Contains(typeValue, "string") {
		toSql := false
		if len(params) == 2 && reflect.TypeOf(params[1]).Kind() == reflect.Bool {
			toSql = params[1].(bool)
		}

		if toSql {
			valueReturn = fmt.Sprintf("'%s'", value)
		} else {
			valueReturn = fmt.Sprintf("%s", value)
		}
	} else if strings.Contains(typeValue, "int") {
		valueReturn = fmt.Sprintf("%d", value)
	} else if strings.Contains(typeValue, "float") {
		valueReturn = fmt.Sprintf("%f", value)
	} else if strings.Contains(typeValue, "bool") {
		valueReturn = fmt.Sprintf("%t", value)
	} else {
		valueReturn = fmt.Sprintf("%v", value)
	}
	return valueReturn
}

func SchemaForUpdate(modelo []Fields) []Fields {
	var newModels []Fields
	for _, v := range modelo {
		if v.PrimaryKey {
			newModels = append(newModels, v)
		} else if v.Where {
			newModels = append(newModels, v)
		} else if v.Update {
			newModels = append(newModels, v)
		}
	}

	return newModels
}

func SchemaForDelete(modelo []Fields) []Fields {
	var newModels []Fields
	for _, v := range modelo {
		if v.PrimaryKey {
			newModels = append(newModels, v)
		} else if v.Where {
			newModels = append(newModels, v)
		}
	}

	return newModels
}
