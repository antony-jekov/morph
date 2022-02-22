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
	"strings"
	"sync"
)

// transformational tags
const (
	//TagTrim trims a string value (e.g "trim" - " value " -> "value")
	TagTrim = "trim"
	//TagLower transforms a string to lower characters (e.g "lower" - "VALUE" -> "value")
	TagLower = "lower"
	//TagUpper transforms a string to upper characters (e.g "upper" - "value" -> "VALUE")
	TagUpper = "upper"
	//TagTruncate truncates a string to a specified length (e.g "truncate=3" - "value" -> "val")
	TagTruncate = "truncate"
	//TagCeil performs ceiling on a floating number (e.g "ceil" - "1.45" -> "2.00")
	TagCeil = "ceil"
	//TagFloor performs flooring on a floating number (e.g "floor" - "1.65" -> "1.00")
	TagFloor = "floor"
	//TagRound performs rounding on a floating number (e.g "round" - "1.45" -> "1.00")
	TagRound = "round"
	//TagPrecision limits precision for a floating number (e.g "precision=2" - "1.499" -> "1.49")
	TagPrecision = "precision"
)

// navigational tags
const (
	//TagDive enters inside slices, arrays or maps to perform transformations on their items, which would've been
	// otherwise neglected - e.g. SomeData []string 'morph:"dive,trim"' - goes inside the array and trims all values
	TagDive = "dive"
	//TagKeys enters keys of a map to perform transformations on them - e.g. SomeData map[string]string
	// 'morph:"dive,keys,trim"' - goes inside the map and dives into its keys to trim them all
	TagKeys = "keys"
	//TagExit states ending of the keys' transformations - e.g. SomeData map[string]string
	// 'morph:"dive,keys,trim,exit,trim"' - goes inside the map and dives into its keys to trim them all after which
	// it exits the keys and trims the values
	TagExit = "exit"
	//TagIgnore ignores a field of type struct and doesn't perform its underlying transformations - e.g. SomeData
	// SomeStruct 'morph:"-"' - ignores this field and doesn't perform its internal morphing
	TagIgnore = "-"
)

const (
	//DefaultTag is the tag used for morphing fields if no other tag is specified.
	DefaultTag = "morph"
	//TagSeparator is the rune for separating the provided tags.
	TagSeparator = ','
	//ParamsSign is the rune that indicates parameters if the tag supports them.
	ParamsSign = '='
)

var navigationalTags = map[string]bool{
	TagDive:   true,
	TagKeys:   true,
	TagExit:   true,
	TagIgnore: true,
}

// Morph transforms the data of a given struct according to a set of provided tags
type Morph interface {

	//Struct accepts a pointer to a struct and performs a transformation on all of its exported fields by morphing
	//them using the provided tags.
	//
	//	Transformational tags:
	//		'trim'      - TagTrim
	//		'lower'     - TagLower
	//		'upper'     - TagUpper
	//		'truncate'  - TagTruncate
	//		'ceil'      - TagCeil
	//		'floor'     - TagFloor
	//		'round'     - TagRound
	//		'precision' - TagPrecision
	//
	//	Navigational tags:
	//		'-'        - TagIgnore
	//		'dive'     - TagDive
	//		'keys'     - TagKeys
	//		'exit'     - TagExit
	//
	//	An example would be:
	//
	//	type EmbeddedModel struct {
	//		EmbeddedField string `morph:"trim"`
	//	}
	//
	//	type InnerModel struct {
	//		SomeOtherField string `morph:"trim"`
	//	}
	//
	//	type Model struct {
	//		EmbeddedModel
	//		SomeField    string `morph:"trim,lower,truncate=5"`
	//		Inner  		 InnerModel
	//		IgnoredInner InnerModel `morph:"-"`
	//		InnerModels  []InnerModel `morph:"dive"`
	//		SomeFields   []string `morph:"dive,trim"`
	//		Numbers      []float64 `morph:"dive,precision=2"`
	//		OtherNumbers []float64 `morph:"dive,floor"`
	//		SomeMap		 map[string]string `morph:"dive,keys,trim,exit,trim"`
	//		SomeOtherMap map[string]InnerModel `morph:"dive,keys,trim,exit"`
	//	}
	//
	//	data := Model {} // fill values
	//	transform := New()
	//	transform.Struct(&data)
	//
	//	Error will be returned if anything else than a pointer to a struct is being passed.
	Struct(structPtr interface{}) error

	// Register accepts custom transformational tags or overrides existing ones and associates the provided
	// transformation function with them.
	// Navigational tags are reserved and are not subject of override. In such case an error will be returned.
	//
	//	Example:
	//		type Model struct {
	//			SomeString string `morph:"swap=baba"`
	//		}
	//		data := Model{SomeString = "value"}
	//
	//		morph := New()
	//		morph.Register("swap", func(value *reflection.Value, params *string) error {
	//			value.SetString(*params)
	//			return nil
	//		})
	//		morph.Struct(&data)
	Register(tag string, transformer FieldTransformer) error

	// WithTag changes the default tag set using DefaultTag to the specified tag if it is valid, otherwise it panics.
	// Valid tags are anything but whitespace.
	//
	//	Example:
	//		type Model struct {
	//			SomeString string `change:"trim"`
	//		}
	//
	//		morph := New().WithTag("change")
	WithTag(tag string) Morph
}

// New creates an instance of Morph with default tags (e.g. TagTrim, TagLower..., etc.)
func New() Morph {
	lock := sync.RWMutex{}
	return &morpher{
		&cache{
			DefaultTag,
			map[string]FieldTransformer{
				TagTrim:  new(trimTransformer),
				TagLower: new(toLowerTransformer),
				TagUpper: new(toUpperTransformer),
				TagTruncate: &truncateTransformer{
					NewIntParamsTransformer(&lock),
				},
				TagCeil:  new(ceilTransformer),
				TagFloor: new(floorTransformer),
				TagRound: new(roundTransformer),
				TagPrecision: &precisionTransformer{
					NewIntParamsTransformer(&lock),
				},
			},
			make(map[string]*structCache),
			&lock,
		},
		&lock,
	}
}

func (c *morpher) WithTag(tag string) Morph {
	tag = strings.TrimSpace(tag)
	if len(tag) == 0 {
		panic(newError(ErrInvalidTagName))
	}

	c.cache.tagName = tag
	return c
}

type morpher struct {
	cache *cache
	mutex *sync.RWMutex
}

func (c *morpher) Register(tag string, transformer FieldTransformer) error {
	tag = strings.TrimSpace(tag)
	if len(tag) == 0 {
		return newError(ErrInvalidTagName)
	}

	if _, ok := navigationalTags[tag]; ok {
		return newErrorf(ErrReservedTagOverride, tag)
	}

	if transformer == nil {
		return newError(ErrInvalidTransformer)
	}

	c.mutex.Lock()
	c.cache.transformers[tag] = transformer
	c.mutex.Unlock()

	return nil
}

func (c *morpher) Struct(structPtr interface{}) error {
	dataValue := reflect.ValueOf(structPtr)
	if dataValue.Kind() != reflect.Ptr {
		return newError(ErrNotAPointer)
	}

	if dataValue.IsNil() {
		return newError(ErrNotAStruct)
	}

	dataValue = dataValue.Elem()
	dataType := dataValue.Type()

	if dataType.Kind() != reflect.Struct {
		return newError(ErrNotAStruct)
	}

	return c.morphStruct(&dataValue, dataType)
}

func (c *morpher) morphStruct(structValue *reflect.Value, structType reflect.Type) error {
	strCache, err := c.cache.getStructCache(structValue, &structType)
	if err != nil {
		return err
	}

	for i := 0; i < strCache.fieldsLength; i++ {
		field := *strCache.fields[i]
		if err = c.morphField(structValue.Field(field.index), field.tags); err != nil {
			return err
		}
	}

	return nil
}

func (c *morpher) morphField(fieldValue reflect.Value, tag *tagChainCache) (err error) {
	actualValue := getActualValue(&fieldValue)
	actualKind := actualValue.Kind()

	if actualKind == reflect.Struct {
		return c.morphStruct(actualValue, actualValue.Type())
	}

	newValue := getAssignableValue(actualValue, &actualKind)
	for currentTag := tag; currentTag != nil && err == nil; currentTag = currentTag.next {
		if currentTag.tag == TagDive {
			err = c.dive(actualValue, &actualKind, currentTag.next)
			break
		}

		if currentTag.transformer == nil {
			continue
		}

		err = currentTag.transformer.Transform(newValue, currentTag.paramsKey)
	}

	if err != nil || newValue != actualValue {
		return
	}

	assignValue(actualValue, newValue, &actualKind)
	return
}

func (c *morpher) dive(actualValue *reflect.Value, actualKind *reflect.Kind, tag *tagChainCache) error {
	switch *actualKind {
	case reflect.Slice, reflect.Array:
		return c.morphCollection(actualValue, tag)
	case reflect.Map:
		return c.morphMap(actualValue, tag)
	}

	return newErrorf(ErrInvalidDiveFmt, actualKind.String())
}

func (c *morpher) morphCollection(sliceValue *reflect.Value, tags *tagChainCache) (err error) {
	itemsLength := sliceValue.Len()
	for i := 0; i < itemsLength && err == nil; i++ {
		err = c.morphField(sliceValue.Index(i), tags)
	}

	return
}

func (c *morpher) morphMap(mapValue *reflect.Value, tags *tagChainCache) error {
	shouldMorphKeys := tags != nil && tags.tag == TagKeys && tags.keysChain != nil
	for _, key := range mapValue.MapKeys() {
		morphedValue := reflect.New(mapValue.Type().Elem()).Elem()
		morphedValue.Set(mapValue.MapIndex(key))

		if shouldMorphKeys {
			mapValue.SetMapIndex(key, reflect.Value{}) // removes key to transform it
			if err := c.morphMapKey(&key, tags.keysChain); err != nil {
				return err
			}
			if err := c.morphField(morphedValue, tags.next); err != nil {
				return err
			}

			mapValue.SetMapIndex(key, morphedValue)
			continue
		}

		if err := c.morphField(morphedValue, tags); err != nil {
			return err
		}

		mapValue.SetMapIndex(key, morphedValue)
	}

	return nil
}

func (c *morpher) morphMapKey(key *reflect.Value, tags *tagChainCache) error {
	morphedKey := reflect.New(key.Type()).Elem()
	morphedKey.Set(*key)
	*key = morphedKey

	return c.morphField(*key, tags)
}
