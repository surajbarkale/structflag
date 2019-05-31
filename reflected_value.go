package structflag

import (
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"strconv"
)

// Value adds ability to get description for flag.Value
type Value interface {
	flag.Getter
	Description() string
}

type reflectedValue struct {
	target      reflect.Value
	description string
}

// NewReflectedValue creates a new flag value that converts string into the given
// reflected value. Bool, Int, UInt and Float values are converted using functions
// from strconv package. For String values, input can be either a bare string or a
// valid JSON string. Arrays, maps and structures must be specified using JSON syntax.
func NewReflectedValue(target reflect.Value, description string) Value {
	return &reflectedValue{target, description}
}

// Description returns stored description for this value.
func (thiz *reflectedValue) Description() string {
	return thiz.description
}

// IsBoolFlag returns true if the required value is boolean. This is added for
// compatibility with kingpin library.
func (thiz *reflectedValue) IsBoolFlag() bool {
	return reflect.Indirect(thiz.target).Kind() == reflect.Bool
}

// String returns the value as string. Primitive values are returned
// as naked values. Complex values are returned as JSON strings.
func (thiz *reflectedValue) String() string {
	return encodeString(thiz.target)
}

// Get returns the underlying value
func (thiz *reflectedValue) Get() interface{} {
	return thiz.target.Interface()
}

// Set updates the value by parsing source string. Complex objects are
// parsed as JSON values.
func (thiz *reflectedValue) Set(source string) error {
	return decodeString(source, thiz.target)
}

func encodeString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return ""
		}
		return encodeString(val.Elem())
	case reflect.String, reflect.Bool, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprint(val.Interface())
	default:
		bytes, err := json.Marshal(val.Interface())
		if err != nil {
			panic(fmt.Errorf("can not convert %s value to string %v", val.Kind().String(), err))
		}
		return string(bytes)
	}
}

func decodeString(s string, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Bool:
		res, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		val.SetBool(res)
	case reflect.Float32, reflect.Float64:
		res, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		if val.OverflowFloat(res) {
			return fmt.Errorf("value %v overflows %s", res, val.Kind().String())
		}
		val.SetFloat(res)
	case reflect.String:
		var res string
		// Try to decode as json string first!
		if json.Unmarshal([]byte(s), &res) != nil {
			res = s
		}
		val.SetString(res)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		if val.OverflowInt(res) {
			return fmt.Errorf("value %v overflows %s", res, val.Kind().String())
		}
		val.SetInt(res)
	case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		if val.OverflowUint(res) {
			return fmt.Errorf("value %v overflows %s", res, val.Kind().String())
		}
		val.SetUint(res)
	case reflect.Ptr:
		res := reflect.New(val.Type().Elem())
		err := decodeString(s, reflect.Indirect(res))
		if err != nil {
			return err
		}
		if val.IsNil() {
			val.Set(res)
		} else {
			val.Elem().Set(res.Elem())
		}
	default:
		res := reflect.New(val.Type())
		err := json.Unmarshal([]byte(s), res.Interface())
		if err != nil {
			return err
		}
		val.Set(res.Elem())
	}
	return nil
}
