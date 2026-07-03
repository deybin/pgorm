package pgorm

import (
	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/clause"
	"github.com/deybin/pgorm/internal/core/ports"
	"github.com/deybin/pgorm/internal/core/services"
)

type TypeJoin = clause.TypeJoin
type OperatorWhere = clause.OperatorWhere
type DBPort = ports.DBPort
type ConfigPgxAdapter = adapters.ConfigPgxAdapter
type SqlExecSingles = services.SqlExecSingles

const SchemaId = adapters.SchemaId
const (
	INNER = clause.INNER
	LEFT  = clause.LEFT
	RIGHT = clause.RIGHT
	FULL  = clause.FULL
)
const (
	I           = clause.I
	D           = clause.D
	MY          = clause.MY
	MYI         = clause.MYI
	MN          = clause.MN
	MNI         = clause.MNI
	LIKE        = clause.LIKE
	IN          = clause.IN
	NOT_IN      = clause.NOT_IN
	BETWEEN     = clause.BETWEEN
	NOT_BETWEEN = clause.NOT_BETWEEN
)
