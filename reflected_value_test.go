package structflag_test

import (
	"flag"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/surajbarkale/structflag"
)

func reflectValue(x interface{}) flag.Value {
	return structflag.NewReflectedValue(reflect.ValueOf(x).Elem(), "")
}

func TestGetValueAsString(t *testing.T) {
	type ts struct {
		A string
		N int
		K struct {
			X float32
			Y float32
		}
	}
	var np *int
	var tr = true
	var f = false
	var i int = 124
	var ip = &i
	var i32 int32 = -8413
	var str string = "hjb ubv73 ,svu83 "
	var list = []float32{0.1, .33, .24, 14215125, 235.58e3}
	var s = &ts{
		A: "shvb",
		N: 325,
		K: struct {
			X float32
			Y float32
		}{
			X: 45.5,
			Y: 3.157,
		},
	}
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"nil", &np, ""},
		{"true", &tr, "true"},
		{"false", &f, "false"},
		{"int", &i, "124"},
		{"intptr", &ip, "124"},
		{"int32", &i32, "-8413"},
		{"string", &str, "hjb ubv73 ,svu83 "},
		{"list", &list, "[0.1,0.33,0.24,14215125,235580]"},
		{"struct", &s, `{"A":"shvb","N":325,"K":{"X":45.5,"Y":3.157}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := reflectValue(tt.value).String()
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestSetStringValue(t *testing.T) {
	src := "test-string abc "
	var val string
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, src, val)
	var ptr *string
	require.NoError(t, reflectValue(&ptr).Set(src))
	require.NotNil(t, ptr)
	assert.Equal(t, src, *ptr)
}

func TestSetJSONStringValue(t *testing.T) {
	exp := "mnbc3yu uqhiq q 33u "
	src := `"` + exp + `"`
	var val string
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, exp, val)
}

func TestSetShouldNotReplacePointer(t *testing.T) {
	src := "test-string abc"
	var val string
	ptr := &val
	require.NoError(t, reflectValue(&ptr).Set(src))
	assert.Equal(t, ptr, &val)
	assert.Equal(t, src, *ptr)
}

func TestErrorOnIntOverflow(t *testing.T) {
	var i int8
	assert.Error(t, reflectValue(&i).Set("555"))

	var j uint8
	assert.Error(t, reflectValue(&j).Set("555"))
}

func TestSetMapFromJSON(t *testing.T) {
	exp := map[string]string{
		"a": "x",
		"b": "y",
	}
	src := `{"a": "x", "b": "y"}`
	var val map[string]string
	var ptr *map[string]string
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, val, exp)
	require.NoError(t, reflectValue(&ptr).Set(src))
	require.NotNil(t, ptr)
	assert.Equal(t, exp, *ptr)
}

func TestSetStructFromJSON(t *testing.T) {
	type ts struct {
		X   int
		Y   int
		Str string
	}
	exp := ts{6, 2, "data"}
	src := `{"x": 6, "y": 2, "str": "data"}`
	var val ts
	var ptr *ts
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, val, exp)
	require.NoError(t, reflectValue(&ptr).Set(src))
	require.NotNil(t, ptr)
	assert.Equal(t, exp, *ptr)
}

func TestSetIntSliceFromJSON(t *testing.T) {
	exp := []int{1, 2, 4, 4}
	src := "[1, 2, 4, 4]"
	var val []int
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, val, exp)
}

func TestSetStringSliceFromJSON(t *testing.T) {
	exp := []string{"a", "p", "cd", "b"}
	src := `["a", "p", "cd", "b"]`
	var val []string
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, val, exp)
}

func TestSetFloat32SliceFromJSON(t *testing.T) {
	exp := []float32{0.5, 0.9, 123, 1e-7}
	src := "[0.5, 0.9, 123, 1e-7]"
	var val []float32
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, val, exp)
}

func TestSetSliceOfSliceFromJSON(t *testing.T) {
	exp := [][]string{{"a", "b"}, {"x"}, {"p", "q", "r"}}
	src := `[["a", "b"], ["x"], ["p", "q", "r"]]`
	var val [][]string
	require.NoError(t, reflectValue(&val).Set(src))
	assert.Equal(t, val, exp)
}

func TestSetReturnsErrorOnInvalidJSON(t *testing.T) {
	src := "[0.5, 0.9, 123, 1e-7"
	var val []float32
	assert.Error(t, reflectValue(&val).Set(src))
}

func TestSetBoolValueWithStrconv(t *testing.T) {
	tests := []string{
		"true",
		"TRUE",
		"t",
		"T",
		"1",
		"false",
		"FALSE",
		"f",
		"F",
		"0",
		"yes",
		"no",
	}
	for _, src := range tests {
		t.Run("bool "+src, func(t *testing.T) {
			exp, expErr := strconv.ParseBool(src)
			var val bool
			var ptr *bool

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(src))
				require.Error(t, reflectValue(&ptr).Set(src))
			} else {
				require.NoError(t, reflectValue(&val).Set(src))
				assert.Equal(t, exp, val)
				require.NoError(t, reflectValue(&ptr).Set(src))
				require.NotNil(t, ptr)
				assert.Equal(t, exp, *ptr)
			}
		})
	}
}

func TestSetIntValueWithStrconv(t *testing.T) {
	tests := []string{
		"invalid",
		"0",
		"+0",
		"-0",
		"1343",
		"+961343",
		"-542253",
		"7553432235352",
		"invalid",
	}
	for _, src := range tests {
		t.Run("int "+src, func(t *testing.T) {
			exp, expErr := strconv.ParseInt(src, 10, 64)
			var val int
			var ptr *int

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(src))
				require.Error(t, reflectValue(&ptr).Set(src))
			} else {
				require.NoError(t, reflectValue(&val).Set(src))
				assert.Equal(t, int(exp), val)
				require.NoError(t, reflectValue(&ptr).Set(src))
				require.NotNil(t, ptr)
				assert.Equal(t, int(exp), *ptr)
			}
		})
	}
	for _, src := range tests {
		t.Run("int32 "+src, func(t *testing.T) {
			exp, expErr := strconv.ParseInt(src, 10, 32)
			var val int32
			var ptr *int32

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(src))
				require.Error(t, reflectValue(&ptr).Set(src))
			} else {
				require.NoError(t, reflectValue(&val).Set(src))
				assert.Equal(t, int32(exp), val)
				require.NoError(t, reflectValue(&ptr).Set(src))
				require.NotNil(t, ptr)
				assert.Equal(t, int32(exp), *ptr)
			}
		})
	}
}

func TestSetUintValueWithStrconv(t *testing.T) {
	tests := []string{
		"0",
		"1343",
		"7553432235352",
		"+1",
		"-1",
		"invalid",
	}
	for _, src := range tests {
		t.Run("uint "+src, func(t *testing.T) {
			exp, expErr := strconv.ParseUint(src, 10, 64)
			var val uint
			var ptr *uint

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(src))
				require.Error(t, reflectValue(&ptr).Set(src))
			} else {
				require.NoError(t, reflectValue(&val).Set(src))
				assert.Equal(t, uint(exp), val)
				require.NoError(t, reflectValue(&ptr).Set(src))
				require.NotNil(t, ptr)
				assert.Equal(t, uint(exp), *ptr)
			}
		})
	}
	for _, src := range tests {
		t.Run("uint32 "+src, func(t *testing.T) {
			exp, expErr := strconv.ParseUint(src, 10, 32)
			var val uint32
			var ptr *uint32

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(src))
				require.Error(t, reflectValue(&ptr).Set(src))
			} else {
				require.NoError(t, reflectValue(&val).Set(src))
				assert.Equal(t, uint32(exp), val)
				require.NoError(t, reflectValue(&ptr).Set(src))
				require.NotNil(t, ptr)
				assert.Equal(t, uint32(exp), *ptr)
			}
		})
	}
}

func TestSetFloat32ValueWithStrconv(t *testing.T) {
	type entry struct {
		name string
		src  string
	}
	tests := []entry{
		{"zero", "0"},
		{"inf", "inf"},
		{"positive", "1343"},
		{"exponent", "325.687e32"},
		{"large", "7553432235352e-156"},
		{"invalid", "invalid"},
		{"overflow", "1e50"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, expErr := strconv.ParseFloat(tt.src, 32)
			var val float32
			var ptr *float32

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(tt.src))
				require.Error(t, reflectValue(&ptr).Set(tt.src))
			} else {
				require.NoError(t, reflectValue(&val).Set(tt.src))
				assert.Equal(t, float32(exp), val)
				require.NoError(t, reflectValue(&ptr).Set(tt.src))
				require.NotNil(t, ptr)
				assert.Equal(t, float32(exp), *ptr)
			}
		})
	}
}

func TestSetFloat64ValueWithStrconv(t *testing.T) {
	type entry struct {
		name string
		src  string
	}
	tests := []entry{
		{"zero", "0"},
		{"inf", "inf"},
		{"positive", "1343"},
		{"exponent", "325.687e32"},
		{"large", "7553432235352e-156"},
		{"invalid", "invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, expErr := strconv.ParseFloat(tt.src, 64)
			var val float64
			var ptr *float64

			if expErr != nil {
				require.Error(t, reflectValue(&val).Set(tt.src))
				require.Error(t, reflectValue(&ptr).Set(tt.src))
			} else {
				require.NoError(t, reflectValue(&val).Set(tt.src))
				assert.Equal(t, exp, val)
				require.NoError(t, reflectValue(&ptr).Set(tt.src))
				require.NotNil(t, ptr)
				assert.Equal(t, exp, *ptr)
			}
		})
	}
}
