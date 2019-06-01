package structflag

import (
	"reflect"
)

// StructToFlagsConverter is useful for converting all fields in a struct to
// a flat map containing flag package compatible values.
type StructToFlagsConverter struct {
	// WordSeparator is used to separate child struct fields from parent structs.
	WordSeparator string
	// DescriptionTag is used to query struct tag to generate description for values.
	DescriptionTag string
	// NameConverterFunc is used to change field names before adding them to output.
	NameConverterFunc func(string) string
}

/*
NewStructToFlagsConverter returns a new converter that uses "-" for separating words,
does not change field names and extracts description from "description" struct tag. The
returned instance can be customized by changing fields. It can be used with flags
package like this:
	package main

	import (
		"flag"

		"github.com/surajbarkale/structflag"
	)

	type extra struct {
		WrapLines bool
		Pages     []int
	}
	type args struct {
		Debug     bool    `description:"Enable debug mode"`
		InputFile *string `description:"Name of input file"`
		Extra     *extra
	}

	func main() {
		a := &args{Debug: true}
		for name, value := range structflag.DefaultStructToFlagsConverter.Convert(&a) {
			flag.Var(value, name, value.Description())
		}
		flag.PrintDefaults()
	}

This program should print output:
	-Debug
		Enable debug mode (default true)
	-Extra-Pages

	-Extra-WrapLines
			(default false)
	-InputFile
		Name of input file
*/
func NewStructToFlagsConverter() *StructToFlagsConverter {
	return &StructToFlagsConverter{
		WordSeparator:     "-",
		DescriptionTag:    "description",
		NameConverterFunc: func(s string) string { return s },
	}
}

// Convert generates the flag values compatible with the structure. You must pass a
// pointer to the value
func (thiz *StructToFlagsConverter) Convert(input interface{}) map[string]Value {
	output := map[string]Value{}
	thiz.reflectStructToFlags("", reflect.ValueOf(input), output)
	return output
}

func (thiz *StructToFlagsConverter) reflectStructToFlags(prefix string, input reflect.Value, output map[string]Value) {
	for input.Kind() == reflect.Ptr || input.Kind() == reflect.Interface {
		input = input.Elem()
	}
	inputType := input.Type()
	for i := 0; i < input.NumField(); i++ {
		field := input.Field(i)
		// Ignore fields that can not be set (i.e. private fields)
		if !field.CanSet() {
			continue
		}
		fieldKind := field.Kind()
		fieldPath := prefix + thiz.NameConverterFunc(inputType.Field(i).Name)
		// Recursively go through the members that are structs or pointers to struct
		if fieldKind == reflect.Struct || (fieldKind == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct) {
			// If struct pointer is nil, then initialize it with empty struct
			if fieldKind == reflect.Ptr && field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			thiz.reflectStructToFlags(fieldPath+thiz.WordSeparator, field, output)
		} else {
			var description string
			if thiz.DescriptionTag != "" {
				description = inputType.Field(i).Tag.Get(thiz.DescriptionTag)
			}
			output[fieldPath] = NewReflectedValue(field, description)
		}
	}
}
