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
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

//FieldTransformer is the actual transformer being called for the fields with a corresponding tag
type FieldTransformer interface {
	//Transform is transforming the given value. The value is a copy of the original in cases where this is needed and
	// is being updated after all the transformations are done.
	//
	// value is the value to be transformed
	//
	// paramsKey is the key for the cached parameters
	Transform(value *reflect.Value, paramsKey *string) error

	//Cache is performing the necessary parsing and converting of the transformer's parameters before they can be used
	//
	// params are the actual parameters
	//
	// paramsKey is the key where the parameters will be stored
	Cache(params, paramsKey *string) error
}

//IntParameterTransformer is used to convert int params and store them for use in the transformation process
type IntParameterTransformer struct {
	Values map[string]*int
	Mutex  *sync.RWMutex
}

//NewIntParamsTransformer returns a new instance
func NewIntParamsTransformer(mutex *sync.RWMutex) IntParameterTransformer {
	return IntParameterTransformer{
		make(map[string]*int),
		mutex,
	}
}

func (t *IntParameterTransformer) Cache(params, key *string) error {
	value, err := strconv.Atoi(*params)
	if err != nil {
		return newErrorf(ErrInvalidParameters, TagPrecision, *params)
	}

	t.Mutex.Lock()
	t.Values[*key] = &value
	t.Mutex.Unlock()

	return nil
}

type ParameterlessTransformer struct {
}

func (t *ParameterlessTransformer) Cache(_, _ *string) error {
	return nil
}

//region Trim

type trimTransformer struct {
	ParameterlessTransformer
}

func (t *trimTransformer) Transform(value *reflect.Value, _ *string) error {
	if value.Kind() != reflect.String {
		return newErrorf(ErrUnexpectedValue, TagTrim, value.Type().Kind().String())
	}

	value.SetString(strings.TrimSpace(value.String()))
	return nil
}

//endregion Trim

//region ToLower

type toLowerTransformer struct {
	ParameterlessTransformer
}

func (t *toLowerTransformer) Transform(value *reflect.Value, _ *string) error {
	if value.Kind() != reflect.String {
		return newErrorf(ErrUnexpectedValue, TagLower, value.Type().Kind().String())
	}

	value.SetString(strings.ToLower(value.String()))
	return nil
}

//endregion ToLower

// region ToUpper

type toUpperTransformer struct {
	ParameterlessTransformer
}

func (t *toUpperTransformer) Transform(value *reflect.Value, _ *string) error {
	if value.Kind() != reflect.String {
		return newErrorf(ErrUnexpectedValue, TagUpper, value.Type().Kind().String())
	}

	value.SetString(strings.ToUpper(value.String()))
	return nil
}

//endregion ToUpper

//region Truncate

type truncateTransformer struct {
	IntParameterTransformer
}

func (t *truncateTransformer) Transform(value *reflect.Value, paramsKey *string) error {
	if value.Kind() != reflect.String {
		return newErrorf(ErrUnexpectedValue, TagTruncate, value.Type().Kind().String())
	}

	t.Mutex.RLock()
	limit, ok := t.Values[*paramsKey]
	t.Mutex.RUnlock()

	if !ok || limit == nil {
		return newErrorf(ErrMissingParametersFmt, TagTruncate)
	}

	if *limit < 0 {
		return newErrorf(ErrInvalidParameters)
	}

	if len(value.String()) > *limit {
		value.SetString(value.String()[:*limit])
	}

	return nil
}

//endregion Truncate

//region Ceil

type ceilTransformer struct {
	ParameterlessTransformer
}

func (t *ceilTransformer) Transform(value *reflect.Value, _ *string) error {
	if value.Kind() != reflect.Float64 && value.Kind() != reflect.Float32 {
		return newErrorf(ErrUnexpectedValue, TagCeil, value.Type().Kind().String())
	}

	value.SetFloat(math.Ceil(value.Float()))

	return nil
}

//endregion Ceil

//region Floor

type floorTransformer struct {
	ParameterlessTransformer
}

func (t *floorTransformer) Transform(value *reflect.Value, _ *string) error {
	if value.Kind() != reflect.Float64 && value.Kind() != reflect.Float32 {
		return newErrorf(ErrUnexpectedValue, TagFloor, value.Type().Kind().String())
	}

	value.SetFloat(math.Floor(value.Float()))

	return nil
}

//endregion Floor

//region Round

type roundTransformer struct {
	ParameterlessTransformer
}

func (t *roundTransformer) Transform(value *reflect.Value, _ *string) error {
	if value.Kind() != reflect.Float64 && value.Kind() != reflect.Float32 {
		return newErrorf(ErrUnexpectedValue, TagRound, value.Type().Kind().String())
	}

	value.SetFloat(math.Round(value.Float()))

	return nil
}

//endregion Round

//region Precision

type precisionTransformer struct {
	IntParameterTransformer
}

func (t *precisionTransformer) Transform(value *reflect.Value, key *string) error {
	if value.Kind() != reflect.Float64 && value.Kind() != reflect.Float32 {
		return newErrorf(ErrUnexpectedValue, TagPrecision, value.Type().Kind().String())
	}

	t.Mutex.RLock()
	precision, ok := t.Values[*key]
	t.Mutex.RUnlock()

	if !ok {
		return newErrorf(ErrMissingParametersFmt, TagPrecision)
	}

	precisionValue := 1.0
	for p := *precision; p > 0; p-- {
		precisionValue *= 10
	}

	value.SetFloat(float64(int(value.Float()*precisionValue)) / precisionValue)

	return nil
}

//endregion Precision
