package adapters

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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
