/******************************************************************************/
/* enc_helpers.go                                                             */
/******************************************************************************/
/*                           This file is part of:                            */
/*                                KAIJU ENGINE                                */
/*                          https://kaijuengine.org                           */
/******************************************************************************/
/* MIT License                                                                */
/*                                                                            */
/* Copyright (c) 2023-present Kaiju Engine authors (AUTHORS.md).              */
/* Copyright (c) 2015-present Brent Farris.                                   */
/*                                                                            */
/* May all those that this source may reach be blessed by the LORD and find   */
/* peace and joy in life.                                                     */
/* Everyone who drinks of this water will be thirsty again; but whoever       */
/* drinks of the water that I will give him shall never thirst; John 4:13-14  */
/*                                                                            */
/* Permission is hereby granted, free of charge, to any person obtaining a    */
/* copy of this software and associated documentation files (the "Software"), */
/* to deal in the Software without restriction, including without limitation  */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,   */
/* and/or sell copies of the Software, and to permit persons to whom the      */
/* Software is furnished to do so, subject to the following conditions:       */
/*                                                                            */
/* The above copyright, blessing, biblical verse, notice and                  */
/* this permission notice shall be included in all copies or                  */
/* substantial portions of the Software.                                      */
/*                                                                            */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS    */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                 */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.     */
/* IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY       */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT  */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE      */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                              */
/******************************************************************************/

// Code generated by go run encgen.go -output enc_helpers.go; DO NOT EDIT.

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gob

import (
	"reflect"
)

var encArrayHelper = map[reflect.Kind]encHelper{
	reflect.Bool:       encBoolArray,
	reflect.Complex64:  encComplex64Array,
	reflect.Complex128: encComplex128Array,
	reflect.Float32:    encFloat32Array,
	reflect.Float64:    encFloat64Array,
	reflect.Int:        encIntArray,
	reflect.Int16:      encInt16Array,
	reflect.Int32:      encInt32Array,
	reflect.Int64:      encInt64Array,
	reflect.Int8:       encInt8Array,
	reflect.String:     encStringArray,
	reflect.Uint:       encUintArray,
	reflect.Uint16:     encUint16Array,
	reflect.Uint32:     encUint32Array,
	reflect.Uint64:     encUint64Array,
	reflect.Uintptr:    encUintptrArray,
}

var encSliceHelper = map[reflect.Kind]encHelper{
	reflect.Bool:       encBoolSlice,
	reflect.Complex64:  encComplex64Slice,
	reflect.Complex128: encComplex128Slice,
	reflect.Float32:    encFloat32Slice,
	reflect.Float64:    encFloat64Slice,
	reflect.Int:        encIntSlice,
	reflect.Int16:      encInt16Slice,
	reflect.Int32:      encInt32Slice,
	reflect.Int64:      encInt64Slice,
	reflect.Int8:       encInt8Slice,
	reflect.String:     encStringSlice,
	reflect.Uint:       encUintSlice,
	reflect.Uint16:     encUint16Slice,
	reflect.Uint32:     encUint32Slice,
	reflect.Uint64:     encUint64Slice,
	reflect.Uintptr:    encUintptrSlice,
}

func encBoolArray(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encBoolSlice(state, v.Slice(0, v.Len()))
}

func encBoolSlice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]bool)
	if !ok {
		// It is kind bool but not type bool. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != false || state.sendZero {
			if x {
				state.encodeUint(1)
			} else {
				state.encodeUint(0)
			}
		}
	}
	return true
}

func encComplex64Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encComplex64Slice(state, v.Slice(0, v.Len()))
}

func encComplex64Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]complex64)
	if !ok {
		// It is kind complex64 but not type complex64. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0+0i || state.sendZero {
			rpart := floatBits(float64(real(x)))
			ipart := floatBits(float64(imag(x)))
			state.encodeUint(rpart)
			state.encodeUint(ipart)
		}
	}
	return true
}

func encComplex128Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encComplex128Slice(state, v.Slice(0, v.Len()))
}

func encComplex128Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]complex128)
	if !ok {
		// It is kind complex128 but not type complex128. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0+0i || state.sendZero {
			rpart := floatBits(real(x))
			ipart := floatBits(imag(x))
			state.encodeUint(rpart)
			state.encodeUint(ipart)
		}
	}
	return true
}

func encFloat32Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encFloat32Slice(state, v.Slice(0, v.Len()))
}

func encFloat32Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]float32)
	if !ok {
		// It is kind float32 but not type float32. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			bits := floatBits(float64(x))
			state.encodeUint(bits)
		}
	}
	return true
}

func encFloat64Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encFloat64Slice(state, v.Slice(0, v.Len()))
}

func encFloat64Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]float64)
	if !ok {
		// It is kind float64 but not type float64. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			bits := floatBits(x)
			state.encodeUint(bits)
		}
	}
	return true
}

func encIntArray(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encIntSlice(state, v.Slice(0, v.Len()))
}

func encIntSlice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]int)
	if !ok {
		// It is kind int but not type int. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeInt(int64(x))
		}
	}
	return true
}

func encInt16Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encInt16Slice(state, v.Slice(0, v.Len()))
}

func encInt16Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]int16)
	if !ok {
		// It is kind int16 but not type int16. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeInt(int64(x))
		}
	}
	return true
}

func encInt32Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encInt32Slice(state, v.Slice(0, v.Len()))
}

func encInt32Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]int32)
	if !ok {
		// It is kind int32 but not type int32. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeInt(int64(x))
		}
	}
	return true
}

func encInt64Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encInt64Slice(state, v.Slice(0, v.Len()))
}

func encInt64Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]int64)
	if !ok {
		// It is kind int64 but not type int64. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeInt(x)
		}
	}
	return true
}

func encInt8Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encInt8Slice(state, v.Slice(0, v.Len()))
}

func encInt8Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]int8)
	if !ok {
		// It is kind int8 but not type int8. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeInt(int64(x))
		}
	}
	return true
}

func encStringArray(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encStringSlice(state, v.Slice(0, v.Len()))
}

func encStringSlice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]string)
	if !ok {
		// It is kind string but not type string. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != "" || state.sendZero {
			state.encodeUint(uint64(len(x)))
			state.b.WriteString(x)
		}
	}
	return true
}

func encUintArray(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encUintSlice(state, v.Slice(0, v.Len()))
}

func encUintSlice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]uint)
	if !ok {
		// It is kind uint but not type uint. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeUint(uint64(x))
		}
	}
	return true
}

func encUint16Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encUint16Slice(state, v.Slice(0, v.Len()))
}

func encUint16Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]uint16)
	if !ok {
		// It is kind uint16 but not type uint16. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeUint(uint64(x))
		}
	}
	return true
}

func encUint32Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encUint32Slice(state, v.Slice(0, v.Len()))
}

func encUint32Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]uint32)
	if !ok {
		// It is kind uint32 but not type uint32. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeUint(uint64(x))
		}
	}
	return true
}

func encUint64Array(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encUint64Slice(state, v.Slice(0, v.Len()))
}

func encUint64Slice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]uint64)
	if !ok {
		// It is kind uint64 but not type uint64. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeUint(x)
		}
	}
	return true
}

func encUintptrArray(state *encoderState, v reflect.Value) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return encUintptrSlice(state, v.Slice(0, v.Len()))
}

func encUintptrSlice(state *encoderState, v reflect.Value) bool {
	slice, ok := v.Interface().([]uintptr)
	if !ok {
		// It is kind uintptr but not type uintptr. TODO: We can handle this unsafely.
		return false
	}
	for _, x := range slice {
		if x != 0 || state.sendZero {
			state.encodeUint(uint64(x))
		}
	}
	return true
}
