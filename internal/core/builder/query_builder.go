package builder

import (
	"strings"

	"github.com/deybin/pgorm/internal/core/domain"
)

func BuildQuery(q *domain.Sintaxis) string {
	var querySql strings.Builder
	// var queryString string
	if !q.WorkQueryFull_field {
		querySql.WriteString(q.Select_field.Build())
		querySql.WriteString(q.From_field.Build())
		querySql.WriteString(q.Join_field.Build())
		/** aplicando Where : where ,and ,or ,in, between ,not in ,not between*/
		querySql.WriteString(q.Where_field.Build())
		q.Args_field = q.Where_field.FindArguments()
		q.ArgsLen_field = q.Where_field.FindArgumentsLen()

		/** aplicando Group by*/
		querySql.WriteString(q.GroupBy_field.Build())

		/** aplicando order by  */
		querySql.WriteString(q.OrderBy_field.Build())
		/** aplicando Top y LImit  */
		querySql.WriteString(q.Limit_field.Build())
	} else {
		querySql.WriteString(q.QueryFull_field)
	}

	return querySql.String()
}
