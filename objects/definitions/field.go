package definitions

type Type string

const (
	TypeInput    Type = "input"
	TypeTextarea Type = "textarea"
	TypeWYSIWYG  Type = "wysiwyg"
	TypePassword Type = "password"

	TypeNumeric Type = "numeric"

	TypeDate     Type = "date"
	TypeDateTime Type = "datetime"
	TypeTime     Type = "time"

	TypeSelect      Type = "select"
	TypeCountry     Type = "country"
	TypeLanguage    Type = "language"
	TypeMultiselect Type = "multiselect"
	TypeCountries   Type = "countries"
	TypeLanguages   Type = "languages"
)

type DBType string

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

type Field struct {
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	DBType      DBType `json:"dbtype"`
	Mandatory   bool   `json:"mandatory"`
	NotEditable bool   `json:"noteditable"`
	Index       bool   `json:"index"`
	Unique      bool   `json:"unique "`
}
