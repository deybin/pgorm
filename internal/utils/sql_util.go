package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func QueryCrossUpdate(query string) string {

	query_regEx := regexp.MustCompile(`ADD_(.*?)_SUMA=`)
	array := query_regEx.FindAllStringSubmatch(query, -1)
	for _, v := range array {

		new := fmt.Sprintf("%s=%s+", v[1], v[1])

		query = strings.Replace(query, v[0], new, -1)
	}

	return query
}
