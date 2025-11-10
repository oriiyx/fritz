package definitions

import "fmt"

type DBType string

// WithSize creates a sized type like varchar(255)
func (d DBType) WithSize(size int) DBType {
	return DBType(fmt.Sprintf("%s(%d)", d, size))
}

// WithPrecision creates types like decimal(10,2)
func (d DBType) WithPrecision(precision, scale int) DBType {
	return DBType(fmt.Sprintf("%s(%d,%d)", d, precision, scale))
}

const (
	DataTypeVarchar DBType = "varchar"
	DataTypeText    DBType = "text"
	DataTypeChar    DBType = "char"

	DataTypeSmallInt DBType = "smallint"
	DataTypeInteger  DBType = "integer"
	DataTypeBigInt   DBType = "bigint"

	DataTypeNumeric DBType = "numeric"
	DataTypeDecimal DBType = "decimal"

	DataTypeFloat4 DBType = "float4"
	DataTypeFloat8 DBType = "float8"

	DataTypeSmallSerial DBType = "smallserial"
	DataTypeSerial      DBType = "serial"
	DataTypeBigSerial   DBType = "bigserial"

	DataTypeBytea DBType = "bytea"
	DataTypeBit   DBType = "bit"

	DataTypeDate        DBType = "date"
	DataTypeTimestamp   DBType = "timestamp"
	DataTypeTimestampTZ DBType = "timestamptz"
	DataTypeTime        DBType = "time"
	DataTypeTimeTZ      DBType = "timetz"
	DataTypeInterval    DBType = "interval"

	DataTypeBoolean DBType = "boolean"
)
