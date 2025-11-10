package definitions

type FieldType string

const (
	TypeInput    FieldType = "input"
	TypeTextarea FieldType = "textarea"
	TypeWYSIWYG  FieldType = "wysiwyg"
	TypePassword FieldType = "password"

	TypeNumeric FieldType = "numeric"

	TypeDate     FieldType = "date"
	TypeDateTime FieldType = "datetime"
	TypeTime     FieldType = "time"

	TypeSelect      FieldType = "select"
	TypeCountry     FieldType = "country"
	TypeLanguage    FieldType = "language"
	TypeMultiselect FieldType = "multiselect"
	TypeCountries   FieldType = "countries"
	TypeLanguages   FieldType = "languages"
)
