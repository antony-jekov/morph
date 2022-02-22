/*
	MIT License

	Copyright (c) 2022 Antony Jekov

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package morph

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

//region Struct

//region bad values
func Test_NilPointer(t *testing.T) {
	transformer := New()
	err := transformer.Struct(nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "not a pointer")
}

func Test_PointerNilStruct(t *testing.T) {
	type Data struct{}
	var ptr *Data
	transformer := New()
	err := transformer.Struct(ptr)

	require.Error(t, err)
	require.Contains(t, err.Error(), "not a struct")
}

func Test_NotAPointer(t *testing.T) {
	type testData struct {
		String string
	}

	data := testData{
		String: "someString",
	}

	transformer := New()
	err := transformer.Struct(data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "not a pointer")
}

func Test_NotAStruct(t *testing.T) {
	data := "string"

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "not a struct")
}

//endregion bad values

//region nothing

func Test_EmptyStruct(t *testing.T) {
	type testData struct{}
	data := testData{}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
}

func Test_StructWithoutTag(t *testing.T) {
	type testData struct {
		String string
	}
	data := testData{
		String: " data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " data ", data.String)
}

//endregion nothing

//region tags

func Test_StructWithEmptyTag(t *testing.T) {
	type testData struct {
		String string `morph:""`
	}
	data := testData{
		String: " data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " data ", data.String)
}

func Test_StructWithUnknownTag(t *testing.T) {
	type testData struct {
		String string `morph:"baba"`
	}
	data := testData{
		String: " data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown tag")
	require.Contains(t, err.Error(), "baba")
}

func Test_StructWithComma(t *testing.T) {
	type testData struct {
		String string `morph:","`
	}
	data := testData{
		String: " data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " data ", data.String)
}

//endregion tags

//region trim

func Test_StructWithEmptyField(t *testing.T) {
	type testData struct {
		String string `morph:"trim"`
	}
	data := testData{}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "", data.String)
}

func Test_StructWithTagTrim(t *testing.T) {
	type testData struct {
		String string `morph:"trim"`
	}
	data := testData{
		String: " data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", data.String)
}

func Test_StructWithTagTrimPrivate(t *testing.T) {
	type testData struct {
		stringPrivate string `morph:"trim"`
	}
	data := testData{
		stringPrivate: " data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " data ", data.stringPrivate)
}

func Test_StructWithTagTrim_Pointer(t *testing.T) {
	type testData struct {
		String *string `morph:"trim"`
	}

	value := " data "
	data := testData{
		String: &value,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", *data.String)
	require.Equal(t, "data", value)
}

func Test_StructWithTagTrim_FieldWithoutTags(t *testing.T) {
	type testData struct {
		String      string `morph:"trim"`
		OtherString string
	}

	data := testData{
		String:      " data ",
		OtherString: " other data ",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", data.String)
	require.Equal(t, " other data ", data.OtherString)
}

func Test_StructWithTagTrim_InnerField(t *testing.T) {
	type otherData struct {
		OtherString string `morph:"trim"`
	}

	type testData struct {
		OtherData otherData
	}

	data := testData{
		OtherData: otherData{OtherString: " other data "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "other data", data.OtherData.OtherString)
}

func Test_IgnoreStructWithTagTrim_InnerField(t *testing.T) {
	type otherData struct {
		OtherString string `morph:"trim"`
	}

	type testData struct {
		OtherData otherData `morph:"-"`
	}

	data := testData{
		OtherData: otherData{OtherString: " other data "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " other data ", data.OtherData.OtherString)
}

//region arrays

func Test_StructWithTagTrim_Array(t *testing.T) {
	type otherData struct {
		OtherString string `morph:"trim"`
	}

	type testData struct {
		OtherData []otherData `morph:"dive"`
	}

	data := testData{
		OtherData: []otherData{{OtherString: " other data "}},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "other data", data.OtherData[0].OtherString)
}

func Test_StructWithTagTrim_ArrayPointer(t *testing.T) {
	type otherData struct {
		OtherString string `morph:"trim"`
	}

	type testData struct {
		OtherData *[]otherData `morph:"dive"`
	}

	data := testData{
		OtherData: &[]otherData{{OtherString: " other data "}},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "other data", (*data.OtherData)[0].OtherString)
}

func Test_StructWithTagTrim_ArrayOfPointers(t *testing.T) {
	type otherData struct {
		OtherString string `morph:"trim"`
	}

	type testData struct {
		OtherData []*otherData `morph:"dive"`
	}

	data := testData{
		OtherData: []*otherData{{OtherString: " other data "}},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "other data", data.OtherData[0].OtherString)
}

func Test_StructWithTagTrim_ArrayOfStrings(t *testing.T) {
	type testData struct {
		OtherData []string `morph:"dive,trim"`
	}

	data := testData{
		OtherData: []string{" data ", " data2 "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", data.OtherData[0])
	require.Equal(t, "data2", data.OtherData[1])
}

func Test_StructWithTagTrim_ArrayOfArrayOfStringsDoubleDive(t *testing.T) {
	type testData struct {
		OtherData [][]string `morph:"dive,dive,trim"`
	}

	data := testData{
		OtherData: [][]string{{" data ", " data2 "}},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", data.OtherData[0][0])
	require.Equal(t, "data2", data.OtherData[0][1])
}

func Test_StructWithTagTrim_ArrayOfStringPointers(t *testing.T) {
	type testData struct {
		OtherData []*string `morph:"dive,trim"`
	}

	strValue := " data "
	strValue2 := " data2 "

	data := testData{
		OtherData: []*string{&strValue, &strValue2},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", *data.OtherData[0])
	require.Equal(t, "data2", *data.OtherData[1])
}

func Test_StructWithTagTrim_ArrayOfNilStringPointers(t *testing.T) {
	type testData struct {
		OtherData []*string `morph:"dive,trim"`
	}

	data := testData{
		OtherData: []*string{nil, nil},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Nil(t, data.OtherData[0])
	require.Nil(t, data.OtherData[1])
}

func Test_StructWithTagTrim_ArrayOfInterfaces(t *testing.T) {
	type testData struct {
		OtherData []interface{} `morph:"dive,trim"`
	}

	type otherData struct {
		OtherField string `morph:"trim"`
	}

	arr := make([]interface{}, 3)
	arr[0] = 5
	arr[1] = " data "
	arr[2] = otherData{" other data "}

	data := testData{
		OtherData: arr,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected value")
}

func Test_StructWithTagTrim_ArrayOfIntegersWithTrim(t *testing.T) {
	type testData struct {
		Data []int `morph:"dive,trim"`
	}

	data := testData{
		Data: []int{
			5, 4,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected value")
}

func Test_StructWithTagTrim_ArrayOfIntegerPointersWithTrim(t *testing.T) {
	type testData struct {
		Data []*int `morph:"dive,trim"`
	}

	num1 := 5
	num2 := 4

	data := testData{
		Data: []*int{
			&num1, &num2,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected value")
}

func Test_StructWithTagTrim_ArrayWithNilStringPointers(t *testing.T) {
	type testData struct {
		OtherData []*string `morph:"dive,trim"`
	}

	str := " data "

	data := testData{
		OtherData: []*string{nil, &str},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Nil(t, data.OtherData[0])
	require.Equal(t, "data", *data.OtherData[1])
}

func Test_StructWithTagTrim_ArrayOfStringsWithoutDive(t *testing.T) {
	type testData struct {
		OtherData []string `morph:"trim"`
	}

	data := testData{
		OtherData: []string{" data "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected value")
	require.Equal(t, " data ", data.OtherData[0])
}

//endregion arrays

//region maps

//region keys

func Test_Struct_MapOfStringsKeys(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive,keys"`
	}

	data := testData{
		Map: map[string]string{
			" key1 ": " value1 ",
			" key2 ": " value2 ",
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " value1 ", data.Map[" key1 "])
	require.Equal(t, " value2 ", data.Map[" key2 "])
}

func Test_Struct_MapOfStringsKeysExit(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive,keys,exit"`
	}

	data := testData{
		Map: map[string]string{
			" key1 ": " value1 ",
			" key2 ": " value2 ",
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " value1 ", data.Map[" key1 "])
	require.Equal(t, " value2 ", data.Map[" key2 "])
}

func Test_Struct_MapOfStringsKeysTrimExit(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive,keys,trim,exit"`
	}

	data := testData{
		Map: map[string]string{
			" key1 ": " value1 ",
			" key2 ": " value2 ",
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " value1 ", data.Map["key1"])
	require.Equal(t, " value2 ", data.Map["key2"])
}

func Test_Struct_MapOfStringsKeysExitTrim(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive,keys,exit,trim"`
	}

	data := testData{
		Map: map[string]string{
			" key1 ": " value1 ",
			" key2 ": " value2 ",
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value1", data.Map[" key1 "])
	require.Equal(t, "value2", data.Map[" key2 "])
}

//endregion keys

//region values

func Test_StructWithTagTrim_MapOfStringsWithoutDive(t *testing.T) {
	type testData struct {
		Map map[string]string
	}

	data := testData{
		Map: map[string]string{"key": " value "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " value ", data.Map["key"])
}

func Test_Struct_MapOfStringsWithDive(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive"`
	}

	data := testData{
		Map: map[string]string{"key": " value "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " value ", data.Map["key"])
}

func Test_Struct_MapOfStringsWithDiveTrim(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive,trim"`
	}

	data := testData{
		Map: map[string]string{
			"key1": " value ",
			"key2": " value2 ",
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value", data.Map["key1"])
	require.Equal(t, "value2", data.Map["key2"])
}

func Test_Struct_MapOfStructsWithDiveTrim(t *testing.T) {
	type mapData struct {
		String string `morph:"trim"`
	}

	type testData struct {
		Map map[string]mapData `morph:"dive"`
	}

	data := testData{
		Map: map[string]mapData{"key": {String: " value "}},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value", data.Map["key"].String)
}

func Test_Struct_MapOfStringPointersWithDiveTrim(t *testing.T) {
	type testData struct {
		Map map[string]*string `morph:"dive,trim"`
	}

	strValue := " value "
	strValue2 := " value2 "

	data := testData{
		Map: map[string]*string{
			"key1": &strValue,
			"key2": &strValue2,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value", *data.Map["key1"])
	require.Equal(t, "value2", *data.Map["key2"])
}

func Test_Struct_MapOfNilStringPointersWithDiveTrim(t *testing.T) {
	type testData struct {
		Map map[string]*string `morph:"dive,trim"`
	}

	data := testData{
		Map: map[string]*string{
			"key1": nil,
			"key2": nil,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Nil(t, data.Map["key1"])
	require.Nil(t, data.Map["key2"])
}

func Test_Struct_MapWithNilStringPointersWithDiveTrim(t *testing.T) {
	type testData struct {
		Map map[string]*string `morph:"dive,trim"`
	}

	str := " data "

	data := testData{
		Map: map[string]*string{
			"key1": nil,
			"key2": &str,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Nil(t, data.Map["key1"])
	require.Equal(t, "data", *data.Map["key2"])
}

func Test_Struct_MapOfStructPointersWithDiveTrim(t *testing.T) {
	type mapData struct {
		String string `morph:"trim"`
	}

	type testData struct {
		Map map[string]*mapData `morph:"dive"`
	}

	data := testData{
		Map: map[string]*mapData{
			"key1": {" value "},
			"key2": {" value2 "},
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value", (*data.Map["key1"]).String)
	require.Equal(t, "value2", (*data.Map["key2"]).String)
}

func Test_Struct_MapOfNilStructPointersWithDiveTrim(t *testing.T) {
	type mapData struct {
		String string `morph:"trim"`
	}

	type testData struct {
		Map map[string]*mapData `morph:"dive"`
	}

	data := testData{
		Map: map[string]*mapData{
			"key1": nil,
			"key2": nil,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Nil(t, data.Map["key1"])
	require.Nil(t, data.Map["key2"])
}

func Test_Struct_MapWithNilStructPointersWithDiveTrim(t *testing.T) {
	type mapData struct {
		String string `morph:"trim"`
	}

	type testData struct {
		Map map[string]*mapData `morph:"dive"`
	}

	data := testData{
		Map: map[string]*mapData{
			"key1": nil,
			"key2": {" value2 "},
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value2", (*data.Map["key2"]).String)
}

func Test_Struct_MapOfInterfaceWithDive(t *testing.T) {
	type mapData struct {
		String string `morph:"trim"`
	}

	type testData struct {
		Map map[string]interface{} `morph:"dive"`
	}

	val1 := mapData{" value "}
	val3 := " value3 "

	data := testData{
		Map: map[string]interface{}{
			"key1": &val1,
			"key2": " value2 ",
			"key3": &val3,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value", data.Map["key1"].(*mapData).String)
}

//endregion values

//endregion maps

//region EmbeddedField

func Test_StructWithTagTrim_EmbeddedField(t *testing.T) {
	type EmbeddedData struct {
		EmbeddedString string `morph:"trim"`
	}
	type testData struct {
		EmbeddedData
	}

	data := testData{
		EmbeddedData{EmbeddedString: " embedded "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "embedded", data.EmbeddedString)
}

func Test_StructWithTagTrim_EmbeddedPrivateField(t *testing.T) {
	type embeddedData struct {
		EmbeddedString string `morph:"trim"`
	}
	type testData struct {
		embeddedData
	}

	data := testData{
		embeddedData{EmbeddedString: " embedded "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, " embedded ", data.EmbeddedString)
}

//endregion EmbeddedField

//endregion trim

//region upper

func Test_StructWithTagUpper(t *testing.T) {
	type testData struct {
		String string `morph:"upper"`
	}
	data := testData{
		String: "data",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "DATA", data.String)
}

//endregion upper

//region truncate

func Test_StructWithTagTruncateBadParameter(t *testing.T) {
	type testData struct {
		String string `morph:"truncate=baba"`
	}

	data := testData{
		String: "123456",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid parameters")
	require.Contains(t, err.Error(), "baba")
}

func Test_StructWithTagTruncateNegativeParameter(t *testing.T) {
	type testData struct {
		String string `morph:"truncate=-1"`
	}

	data := testData{
		String: "123456",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid parameters")
}

func Test_StructWithTagTruncate(t *testing.T) {
	type testData struct {
		String string `morph:"truncate=5"`
	}
	data := testData{
		String: "123456",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "12345", data.String)
}

func Test_StructWithTagTruncateExact(t *testing.T) {
	type testData struct {
		String string `morph:"truncate=6"`
	}
	data := testData{
		String: "123456",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "123456", data.String)
}

func Test_StructWithTagTruncateOneMoreThanExact(t *testing.T) {
	type testData struct {
		String string `morph:"truncate=7"`
	}
	data := testData{
		String: "123456",
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "123456", data.String)
}

//endregion truncate

//region mixed

func Test_StructWithTagTrimLower_ArrayOfStrings(t *testing.T) {
	type testData struct {
		OtherData []string `morph:"dive,trim,lower"`
	}

	data := testData{
		OtherData: []string{" DATA ", " DATA2 "},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", data.OtherData[0])
	require.Equal(t, "data2", data.OtherData[1])
}

func Test_StructWithTagTrimLower_ArrayOfArrayOfStringsDoubleDive(t *testing.T) {
	type testData struct {
		OtherData [][]string `morph:"dive,dive,trim,lower"`
	}

	data := testData{
		OtherData: [][]string{{" DATA ", " DATA2 "}},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "data", data.OtherData[0][0])
	require.Equal(t, "data2", data.OtherData[0][1])
}

func Test_Struct_MapOfStringsKeysMultipleTransforms(t *testing.T) {
	type testData struct {
		Map map[string]string `morph:"dive,keys,trim,lower,exit,trim,lower"`
	}

	data := testData{
		Map: map[string]string{
			" KEY1 ": " VALUE1 ",
			" KEY2 ": " VALUE2 ",
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "value1", data.Map["key1"])
	require.Equal(t, "value2", data.Map["key2"])
}

//endregion mixed

//region numbers

func Test_Ceil(t *testing.T) {
	type testData struct {
		Numbers64 []float64 `morph:"dive,ceil"`
		Numbers32 []float32 `morph:"dive,ceil"`
	}

	data := testData{
		Numbers64: []float64{
			1.465,
			-1.465,
			1,
		},
		Numbers32: []float32{
			1.465,
			-1.465,
			1,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)

	require.Equal(t, float64(2), data.Numbers64[0])
	require.Equal(t, float64(-1), data.Numbers64[1])
	require.Equal(t, float64(1), data.Numbers64[2])

	require.Equal(t, float32(2), data.Numbers32[0])
	require.Equal(t, float32(-1), data.Numbers32[1])
	require.Equal(t, float32(1), data.Numbers32[2])
}

func Test_Floor(t *testing.T) {
	type testData struct {
		Numbers64 []float64 `morph:"dive,floor"`
		Numbers32 []float32 `morph:"dive,floor"`
	}

	data := testData{
		Numbers64: []float64{
			1.565,
			-1.465,
			1,
		},
		Numbers32: []float32{
			1.465,
			-1.465,
			1,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)

	require.Equal(t, float64(1), data.Numbers64[0])
	require.Equal(t, float64(-2), data.Numbers64[1])
	require.Equal(t, float64(1), data.Numbers64[2])

	require.Equal(t, float32(1), data.Numbers32[0])
	require.Equal(t, float32(-2), data.Numbers32[1])
	require.Equal(t, float32(1), data.Numbers32[2])
}

func Test_Round(t *testing.T) {
	type testData struct {
		Numbers64 []float64 `morph:"dive,round"`
		Numbers32 []float32 `morph:"dive,round"`
	}

	data := testData{
		Numbers64: []float64{
			1.565,
			1.456,
			-1.565,
			-1.465,
			1,
		},
		Numbers32: []float32{
			1.565,
			1.465,
			-1.565,
			-1.465,
			1,
		},
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)

	require.Equal(t, float64(2), data.Numbers64[0])
	require.Equal(t, float64(1), data.Numbers64[1])
	require.Equal(t, float64(-2), data.Numbers64[2])
	require.Equal(t, float64(-1), data.Numbers64[3])
	require.Equal(t, float64(1), data.Numbers64[4])

	require.Equal(t, float32(2), data.Numbers32[0])
	require.Equal(t, float32(1), data.Numbers32[1])
	require.Equal(t, float32(-2), data.Numbers32[2])
	require.Equal(t, float32(-1), data.Numbers32[3])
	require.Equal(t, float32(1), data.Numbers32[4])
}

func Test_PrecisionBadParameter(t *testing.T) {
	type testData struct {
		Num float64 `morph:"precision=baba"`
	}

	data := testData{
		Num: 1.16,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid parameters")
	require.Contains(t, err.Error(), "baba")
}

func Test_PrecisionNoParameter(t *testing.T) {
	type testData struct {
		Num float64 `morph:"precision"`
	}

	data := testData{
		Num: 1.16,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid parameters")
}

func Test_PrecisionMissingParameter(t *testing.T) {
	type testData struct {
		Num float64 `morph:"precision="`
	}

	data := testData{
		Num: 1.16,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid parameters")
}

func Test_Precision(t *testing.T) {
	type testData struct {
		Num164 float64 `morph:"precision=1"`
		Num132 float32 `morph:"precision=1"`
		Num264 float64 `morph:"precision=2"`
		Num232 float32 `morph:"precision=2"`
		Num364 float64 `morph:"precision=0"`
		Num332 float32 `morph:"precision=0"`
	}

	data := testData{
		Num164: 1.16,
		Num132: 1.16,
		Num264: 1.167,
		Num232: 1.167,
		Num364: 1.9,
		Num332: 1.9,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)

	require.Equal(t, 1.1, data.Num164)
	require.Equal(t, 1.16, data.Num264)
	require.Equal(t, float64(1), data.Num364)
	require.Equal(t, float32(1.1), data.Num132)
	require.Equal(t, float32(1.16), data.Num232)
	require.Equal(t, float32(1), data.Num332)
}

func Test_PrecisionOnZeroValue(t *testing.T) {
	type testData struct {
		Num1 float64 `morph:"precision=1"`
		Num2 float32 `morph:"precision=1"`
	}

	data := testData{
		Num1: 0.0,
		Num2: 0.0,
	}

	transformer := New()
	err := transformer.Struct(&data)

	require.Nil(t, err)

	require.Equal(t, 0.0, data.Num1)
	require.Equal(t, float32(0.0), data.Num2)
}

//endregion numbers

//endregion Struct

//region Register

func Test_RegisterDiveOverride(t *testing.T) {
	transformer := New()
	err := transformer.Register("dive", new(emptyTransformer))

	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved tag")
}

func Test_RegisterKeysOverride(t *testing.T) {
	transformer := New()
	err := transformer.Register("keys", new(emptyTransformer))

	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved tag")
}

func Test_RegisterExitKeysOverride(t *testing.T) {
	transformer := New()
	err := transformer.Register("exit", new(emptyTransformer))

	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved tag")
}

func Test_RegisterIgnoreOverride(t *testing.T) {
	transformer := New()
	err := transformer.Register("-", new(emptyTransformer))

	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved tag")
}

func Test_StructWithCustomTag(t *testing.T) {
	type testData struct {
		String string `morph:"baba"`
	}

	data := testData{
		String: " data ",
	}

	transformer := New()
	require.Nil(t, transformer.Register("baba", &funcTransformer{
		Func: func(s *reflect.Value, key *string) error {
			s.SetString("baba")
			return nil
		},
	}))
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "baba", data.String)
}

//endregion Register

//region WithTag

func Test_WithTagEmpty(t *testing.T) {
	require.Panics(t, func() {
		New().WithTag("")
	})
}

func Test_WithTagBlank(t *testing.T) {
	require.Panics(t, func() {
		New().WithTag("     ")
	})
}

func Test_WithTag(t *testing.T) {
	type testData struct {
		SomeString string `change:"upper"`
	}

	data := testData{
		SomeString: "yes",
	}

	transformer := New().WithTag("change")
	err := transformer.Struct(&data)

	require.Nil(t, err)
	require.Equal(t, "YES", data.SomeString)
}

//endregion WithTag

type emptyTransformer struct {
	ParameterlessTransformer
}

type funcTransformer struct {
	ParameterlessTransformer
	Func func(_ *reflect.Value, _ *string) error
}

func (t *funcTransformer) Transform(value *reflect.Value, paramsKey *string) error {
	return t.Func(value, paramsKey)
}

func (t *emptyTransformer) Transform(_ *reflect.Value, _ *string) error {
	return nil
}
