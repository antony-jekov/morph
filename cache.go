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
	"strings"
	"sync"
)

type tagChainCache struct {
	tag         string
	params      *string
	paramsKey   *string
	transformer FieldTransformer
	next        *tagChainCache
	keysChain   *tagChainCache
}

type fieldCache struct {
	index int
	tags  *tagChainCache
}

type structCache struct {
	fieldsLength int
	fields       []*fieldCache
}

type cache struct {
	tagName      string
	transformers map[string]FieldTransformer
	structsCache map[string]*structCache
	mutex        *sync.RWMutex
}

func (c *cache) getStructCache(structValue *reflect.Value, structType *reflect.Type) (*structCache, error) {
	key := fmt.Sprintf("%s/%s", (*structType).PkgPath(), (*structType).Name())

	// safe read
	c.mutex.RLock()
	strCache, ok := c.structsCache[key]
	c.mutex.RUnlock()

	if !ok {
		newCache, err := c.buildStructCache(structValue)
		if err != nil {
			return nil, err
		}

		strCache = newCache

		c.mutex.Lock()
		c.structsCache[key] = strCache
		c.mutex.Unlock()
	}

	return strCache, nil
}

func (c *cache) buildStructCache(structValue *reflect.Value) (*structCache, error) {
	fields := make([]*fieldCache, 0)
	fieldsLength := structValue.NumField()
	strutType := structValue.Type()

	for i := 0; i < fieldsLength; i++ {
		field := strutType.Field(i)

		if !field.IsExported() {
			continue
		}

		tagsRaw := field.Tag.Get(c.tagName)
		if tagsRaw == TagIgnore {
			continue
		}

		var tags *tagChainCache
		if len(tagsRaw) > 0 {
			paramsKey := getParamsKey(structValue.Type().String(), i)
			tagsCache, err := c.buildTagsCache(&tagsRaw, &paramsKey)
			if err != nil {
				return nil, err
			}

			if tagsCache != nil {
				tagsCache.paramsKey = &paramsKey
			}

			tags = tagsCache
		}

		fields = append(fields, &fieldCache{
			index: i,
			tags:  tags,
		})
	}

	return &structCache{
		fieldsLength: len(fields),
		fields:       fields,
	}, nil
}

func (c *cache) buildTagsCache(tagsRaw, paramsKey *string) (*tagChainCache, error) {
	allTags := strings.FieldsFunc(*tagsRaw, func(r rune) bool {
		return r == TagSeparator
	})

	tags := &tagChainCache{}
	currentTag := tags

	for i := 0; i < len(allTags); i++ {
		tag := allTags[i]
		newTagCache, err := c.buildTagCache(tag, paramsKey)
		if err != nil {
			return nil, err
		}

		if tag == TagKeys && i+1 < len(allTags) {
			i++
			keyTagCache := &tagChainCache{}
			currentKeyTagCache := keyTagCache
			for ; i < len(allTags); i++ {
				keyTag := allTags[i]
				if keyTag == TagExit {
					break
				}

				newKeyTagCache, errBuild := c.buildTagCache(keyTag, paramsKey)
				if errBuild != nil {
					return nil, errBuild
				}

				currentKeyTagCache.next = newKeyTagCache
				currentKeyTagCache = newKeyTagCache
			}

			newTagCache.keysChain = keyTagCache.next
		}

		currentTag.next = newTagCache
		currentTag = newTagCache
	}

	return tags.next, nil
}

func (c *cache) buildTagCache(tag string, paramsKey *string) (*tagChainCache, error) {
	params := ""
	equalSignIndex := strings.IndexRune(tag, ParamsSign)

	if equalSignIndex > 0 {
		params = tag[equalSignIndex+1:]
		tag = tag[:equalSignIndex]
	}

	c.mutex.RLock()
	tr, ok := c.transformers[tag]
	c.mutex.RUnlock()

	if !ok && !(navigationalTags[tag]) {
		return nil, newErrorf(ErrUnknownTagFmt, tag)
	}

	if tr != nil {
		if err := tr.Cache(&params, paramsKey); err != nil {
			return nil, err
		}
	}

	return &tagChainCache{
		tag:         tag,
		params:      &params,
		transformer: tr,
	}, nil
}
