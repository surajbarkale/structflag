package structflag_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/surajbarkale/structflag"
)

type nested struct {
	Int      int
	IntPtr   *int
	Float    float32
	FloatPtr *float32
	private  int
}
type param struct {
	Nested    nested
	NestedPtr *nested
	String    string
	StringPtr *string
	IntArray  []int
}

func TestDefaultStructKeys(t *testing.T) {
	val := &param{}
	c := structflag.NewStructToFlagsConverter()
	sv := c.Convert(val)
	expKeys := []string{
		"Nested-Int",
		"Nested-IntPtr",
		"Nested-Float",
		"Nested-FloatPtr",
		"NestedPtr-Int",
		"NestedPtr-IntPtr",
		"NestedPtr-Float",
		"NestedPtr-FloatPtr",
		"String",
		"StringPtr",
		"IntArray",
	}
	assert := assert.New(t)
	assert.Equal(len(expKeys), len(sv))
	for _, k := range expKeys {
		assert.Contains(sv, k)
	}
}

func TestCustomSeparator(t *testing.T) {
	val := &param{}
	c := structflag.NewStructToFlagsConverter()
	c.WordSeparator = "."
	sv := c.Convert(val)
	expKeys := []string{
		"Nested.Int",
		"Nested.IntPtr",
		"Nested.Float",
		"Nested.FloatPtr",
		"NestedPtr.Int",
		"NestedPtr.IntPtr",
		"NestedPtr.Float",
		"NestedPtr.FloatPtr",
		"String",
		"StringPtr",
		"IntArray",
	}
	assert := assert.New(t)
	assert.Equal(len(expKeys), len(sv))
	for _, k := range expKeys {
		assert.Contains(sv, k)
	}
}

func TestCustomNameFunction(t *testing.T) {
	val := &param{}
	c := structflag.NewStructToFlagsConverter()
	c.NameConverterFunc = strings.ToUpper
	sv := c.Convert(val)
	expKeys := []string{
		"NESTED-INT",
		"NESTED-INTPTR",
		"NESTED-FLOAT",
		"NESTED-FLOATPTR",
		"NESTEDPTR-INT",
		"NESTEDPTR-INTPTR",
		"NESTEDPTR-FLOAT",
		"NESTEDPTR-FLOATPTR",
		"STRING",
		"STRINGPTR",
		"INTARRAY",
	}
	assert := assert.New(t)
	assert.Equal(len(expKeys), len(sv))
	for _, k := range expKeys {
		assert.Contains(sv, k)
	}
}

func TestDescriptionTag(t *testing.T) {
	val := struct {
		Param1 int    `usage:"The first parameter"`
		Input  string `usage:"Input string"`
	}{}
	c := structflag.NewStructToFlagsConverter()
	c.DescriptionTag = "usage"
	sv := c.Convert(&val)
	assert := assert.New(t)
	assert.Equal("The first parameter", sv["Param1"].Description())
	assert.Equal("Input string", sv["Input"].Description())
}
