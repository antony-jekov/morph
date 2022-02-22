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
	"fmt"
	"reflect"
)

func getActualValue(dataValue *reflect.Value) *reflect.Value {
	switch dataValue.Kind() {
	case reflect.Ptr, reflect.Interface:
		if dataValue.IsNil() {
			return dataValue
		}

		dataValueElement := dataValue.Elem()
		return getActualValue(&dataValueElement)
	default:
		return dataValue
	}
}

func getAssignableValue(value *reflect.Value, kind *reflect.Kind) *reflect.Value {
	newValue := *value
	if *kind == reflect.Ptr {
		newValue = reflect.New(value.Type().Elem()).Elem()
		if !value.IsNil() {
			newValue.Set(reflect.Indirect(*value))
		}
	} else if !value.CanAddr() {
		newValue = reflect.New(value.Type()).Elem()
		newValue.Set(*value)
	}

	return &newValue
}

func assignValue(target, value *reflect.Value, targetKind *reflect.Kind) {
	if *targetKind == reflect.Ptr {
		(*target).Set(value.Addr())
	} else if !target.CanAddr() {
		target.Set(reflect.Indirect(*value))
	}
}

func getParamsKey(valueType string, fieldIndex int) string {
	return fmt.Sprintf("%s.%d", valueType, fieldIndex)
}
