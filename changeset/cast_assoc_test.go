package changeset

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

type Inner struct {
	Field4 int
	Field5 string
}

func TestCastAssocOne(t *testing.T) {
	var inner Inner
	var entity struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": map[string]interface{}{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, params["field3"].(map[string]interface{})),
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(inner),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": inner,
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssocOnePointer(t *testing.T) {
	var inner Inner
	var entity struct {
		Field1 int
		Field2 string
		Field3 *Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": map[string]interface{}{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, params["field3"].(map[string]interface{})),
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(inner),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssocOneErrorParamsNotAMap(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": "3",
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssocOneInnerChangesetError(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": map[string]interface{}{
			"field4": "4",
		},
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field4 is invalid", ch.Error().Error())
	assert.Equal(t, "field3.field4", ch.Error().(errors.Error).Field)
}

func TestCastAssocMany(t *testing.T) {
	var inner Inner
	var entity struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []map[string]interface{}{
			{
				"field4": 14,
				"field5": "15",
			},
			{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	interfaceParams := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []interface{}{
			map[string]interface{}{
				"field4": 14,
				"field5": "15",
			},
			map[string]interface{}{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	field3 := params["field3"].([]map[string]interface{})
	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []*Changeset{
			changeInner(inner, field3[0]),
			changeInner(inner, field3[1]),
		},
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf([]Inner{}),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": []Inner(nil),
	}

	// with map assoc
	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())

	// with slice of interface assoc
	ch = Cast(entity, interfaceParams, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssocManyPointer(t *testing.T) {
	var inner Inner
	var entity struct {
		Field1 int
		Field2 string
		Field3 []*Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []map[string]interface{}{
			{
				"field4": 14,
				"field5": "15",
			},
			{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	interfaceParams := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []interface{}{
			map[string]interface{}{
				"field4": 14,
				"field5": "15",
			},
			map[string]interface{}{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	field3 := params["field3"].([]map[string]interface{})
	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []*Changeset{
			changeInner(inner, field3[0]),
			changeInner(inner, field3[1]),
		},
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf([]Inner{}),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": []*Inner(nil),
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())

	ch = Cast(entity, interfaceParams, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}
func TestCastAssocManyErrorParamsNotASliceOfAMap(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": "3",
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssocManyErrorMixed(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []interface{}{
			map[string]interface{}{
				"field4": 14,
				"field5": "15",
			},
			"3",
		},
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssocManyInnerChangesetError(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []map[string]interface{}{
			{
				"field4": "14",
			},
		},
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field4 is invalid", ch.Error().Error())
	assert.Equal(t, "field3[0].field4", ch.Error().(errors.Error).Field)
}
