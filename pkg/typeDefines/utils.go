package typeDefines

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

func tryParseTime(s string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02", // "YYYY-MM-DD"
		"02/01/2006", // "DD/MM/YYYY"
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func toFloat64(v any) (float64, bool) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), true
	case reflect.Float32, reflect.Float64:
		return val.Float(), true
	case reflect.String:
		if ret, err := strconv.ParseFloat(val.String(), 64); err == nil {
			return ret, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func ContainsKeyRecursevely(responseVal any, targetVal string) (any, bool) {
	switch v := responseVal.(type) {
	case map[string]any:
		for key, val := range v {
			if key == targetVal {
				return v, true
			}
			if _, ok := ContainsKeyRecursevely(val, targetVal); ok {
				return v, true
			}
		}
	case []any:
		for _, item := range v {
			if _, ok := ContainsKeyRecursevely(item, targetVal); ok {
				return v, true
			}
		}
	}
	return nil, false
}

func StringifyMyAny(myAny any) string {
	switch v := myAny.(type) {
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			log.Printf("error marshalling map: %v", err)
			return ""
		}
		return string(b)
	default:
		return fmt.Sprintf("%v", myAny)
	}
}

func getSpecificVal(fa []string, m any) (any, bool) {
	if len(fa) == 0 {
		return nil, false
	}
	switch val := m.(type) {
	case map[string]any:

		if len(fa) == 1 {
			return val[fa[0]], true
		} else {
			if val, ok := val[fa[0]]; ok {
				return getSpecificVal(fa[1:], val)
			} else {
				return nil, false
			}
		}
	case []any:
		for i := range val {
			if val[i].(map[string]any)[fa[0]] != nil {
				if len(fa) == 1 {
					return (val[i].(map[string]any)[fa[0]]), true
				} else {
					return getSpecificVal(fa[1:], val[i])
				}
			}
		}
	}
	return nil, false
}

func typeOfValue(v any) string {
	switch v.(type) {
	case map[string]any, map[string]string:
		return "object"
	case float64:
		return "number"
	case string:
		_, vFloat := toFloat64(v)
		if vFloat {
			return "number"
		}
		_, vDate := tryParseTime(v.(string))
		if vDate {
			return "date"
		}
		return "string"
	case []float64, []string, []map[string]any, []map[string]string:
		return "array"
	default:
		return "ERROR"
	}
}
