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
)

const (
	ErrNotAPointer        = "the provided value is not a pointer"
	ErrNotAStruct         = "the provided value is not a struct"
	ErrInvalidTagName     = "invalid tag name"
	ErrInvalidTransformer = "invalid transformer"
)

const (
	ErrUnknownTagFmt        = "unknown tag: '%s'"
	ErrInvalidDiveFmt       = "cannot dive into kind: %s"
	ErrUnexpectedValue      = "unexpected value:'%s' for tag: '%s'"
	ErrReservedTagOverride  = "cannot override reserved tag: '%s'"
	ErrInvalidParameters    = "invalid parameters '%s' for tag: '%s'"
	ErrMissingParametersFmt = "missing parameters for tag: %s"
)

type ErrMorph struct {
	message string
}

func (c ErrMorph) Error() string {
	return c.message
}

func newError(message string) error {
	return ErrMorph{message}
}

func newErrorf(messageFmt string, args ...string) error {
	return ErrMorph{fmt.Sprintf(messageFmt, args)}
}
