package pgorm

import (
	"github.com/deybin/pgorm/internal/adapters"
	"github.com/deybin/pgorm/internal/core/clause"
	"github.com/deybin/pgorm/internal/core/ports"
)

type TypeJoin = clause.TypeJoin

const (
	INNER = clause.INNER
	LEFT  = clause.LEFT
	RIGHT = clause.RIGHT
	FULL  = clause.FULL
)

type OperatorWhere = clause.OperatorWhere

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

type DBPort = ports.DBPort

type ConfigPgxAdapter = adapters.ConfigPgxAdapter

const SchemaId = adapters.SchemaId
