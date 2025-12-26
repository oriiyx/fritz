package adapters

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// String helpers - existing ones remain the same
func mustGetString(data map[string]interface{}, key string) string {
	v, ok := data[key]
	if !ok {
		panic(fmt.Sprintf("missing required field: %s", key))
	}
	s, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("field %s must be string, got %T", key, v))
	}
	return s
}

func getString(data map[string]interface{}, key string) string {
	v, ok := data[key]
	if !ok {
		return ""
	}
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

func getPgText(data map[string]interface{}, key string) pgtype.Text {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Text{Valid: false}
	}
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}
	return pgtype.Text{String: s, Valid: true}
}

// Integer helpers - mandatory (NOT NULL)
func mustGetInt32(data map[string]interface{}, key string) int32 {
	v, ok := data[key]
	if !ok {
		panic(fmt.Sprintf("missing required field: %s", key))
	}

	switch val := v.(type) {
	case float64:
		return int32(val)
	case int:
		return int32(val)
	case int32:
		return val
	case int64:
		return int32(val)
	default:
		panic(fmt.Sprintf("field %s must be numeric, got %T", key, v))
	}
}

func mustGetInt64(data map[string]interface{}, key string) int64 {
	v, ok := data[key]
	if !ok {
		panic(fmt.Sprintf("missing required field: %s", key))
	}

	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	default:
		panic(fmt.Sprintf("field %s must be numeric, got %T", key, v))
	}
}

func mustGetInt16(data map[string]interface{}, key string) int16 {
	v, ok := data[key]
	if !ok {
		panic(fmt.Sprintf("missing required field: %s", key))
	}

	switch val := v.(type) {
	case float64:
		return int16(val)
	case int:
		return int16(val)
	case int16:
		return val
	case int32:
		return int16(val)
	case int64:
		return int16(val)
	default:
		panic(fmt.Sprintf("field %s must be numeric, got %T", key, v))
	}
}

// Integer helpers - nullable (pgtype.Int4, pgtype.Int8, pgtype.Int2)
func getPgInt4(data map[string]interface{}, key string) pgtype.Int4 {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Int4{Valid: false}
	}

	var intVal int32
	switch val := v.(type) {
	case float64:
		intVal = int32(val)
	case int:
		intVal = int32(val)
	case int32:
		intVal = val
	case int64:
		intVal = int32(val)
	default:
		return pgtype.Int4{Valid: false}
	}

	return pgtype.Int4{Int32: intVal, Valid: true}
}

func getPgInt8(data map[string]interface{}, key string) pgtype.Int8 {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Int8{Valid: false}
	}

	var intVal int64
	switch val := v.(type) {
	case float64:
		intVal = int64(val)
	case int:
		intVal = int64(val)
	case int32:
		intVal = int64(val)
	case int64:
		intVal = val
	default:
		return pgtype.Int8{Valid: false}
	}

	return pgtype.Int8{Int64: intVal, Valid: true}
}

func getPgInt2(data map[string]interface{}, key string) pgtype.Int2 {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Int2{Valid: false}
	}

	var intVal int16
	switch val := v.(type) {
	case float64:
		intVal = int16(val)
	case int:
		intVal = int16(val)
	case int16:
		intVal = val
	case int32:
		intVal = int16(val)
	case int64:
		intVal = int16(val)
	default:
		return pgtype.Int2{Valid: false}
	}

	return pgtype.Int2{Int16: intVal, Valid: true}
}

// Boolean helpers
func mustGetBool(data map[string]interface{}, key string) bool {
	v, ok := data[key]
	if !ok {
		panic(fmt.Sprintf("missing required field: %s", key))
	}
	b, ok := v.(bool)
	if !ok {
		panic(fmt.Sprintf("field %s must be bool, got %T", key, v))
	}
	return b
}

func getPgBool(data map[string]interface{}, key string) pgtype.Bool {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Bool{Valid: false}
	}

	b, ok := v.(bool)
	if !ok {
		return pgtype.Bool{Valid: false}
	}

	return pgtype.Bool{Bool: b, Valid: true}
}

// Date/Time helpers
func getPgDate(data map[string]interface{}, key string) pgtype.Date {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Date{Valid: false}
	}

	var t time.Time
	switch val := v.(type) {
	case string:
		parsed, err := time.Parse("2006-01-02", val)
		if err != nil {
			return pgtype.Date{Valid: false}
		}
		t = parsed
	case time.Time:
		t = val
	default:
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func getPgTimestamp(data map[string]interface{}, key string) pgtype.Timestamptz {
	v, ok := data[key]
	if !ok || v == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	var t time.Time
	switch val := v.(type) {
	case string:
		parsed, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return pgtype.Timestamptz{Valid: false}
		}
		t = parsed
	case time.Time:
		t = val
	default:
		return pgtype.Timestamptz{Valid: false}
	}

	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
