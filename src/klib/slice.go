/******************************************************************************/
/* slice.go                                                                   */
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

package klib

import (
	"math/rand/v2"
	"slices"
	"unsafe"
)

func RemoveUnordered[T any](slice []T, idx int) []T {
	last := len(slice) - 1
	slice[idx] = slice[last]
	return slice[:last]
}

func Shuffle[T any](slice []T, rng *rand.Rand) {
	if rng == nil {
		for i := len(slice) - 1; i > 0; i-- {
			j := rand.IntN(i + 1)
			slice[i], slice[j] = slice[j], slice[i]
		}
	} else {
		for i := len(slice) - 1; i > 0; i-- {
			j := rng.IntN(i + 1)
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
}

func ShuffleRandom[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func Contains[T comparable](slice []T, item T) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}

func AppendUnique[T comparable](slice []T, values ...T) []T {
	for i := range values {
		if !slices.Contains(slice, values[i]) {
			slice = append(slice, values[i])
		}
	}
	return slice
}

func ByteSliceToFloat32Slice(data []byte) []float32 {
	fLen := len(data) / int(unsafe.Sizeof(float32(0)))
	f := *(*[]float32)(unsafe.Pointer(&data))
	return f[:fLen:fLen]
}

func ByteSliceToUInt16Slice(data []byte) []uint16 {
	ui16Len := len(data) / int(unsafe.Sizeof(uint16(0)))
	u := *(*[]uint16)(unsafe.Pointer(&data))
	return u[:ui16Len:ui16Len]
}

func RemoveNils[S any](slice []*S) []*S {
	result := make([]*S, 0, len(slice))
	for i := range slice {
		if slice[i] != nil {
			result = append(result, slice[i])
		}
	}
	return result
}

func SliceMove[S any](s []S, from, to int) {
	if from == to {
		return
	} else if to == from-1 || to == from+1 {
		s[to], s[from] = s[from], s[to]
	} else if to < from {
		a, b := to, from+1
		temp := s[b-1]
		copy(s[a+1:b], s[a:b-1])
		s[a] = temp
	} else if to > from {
		a, b := from, to+1
		temp := s[a]
		copy(s[a:b-1], s[a+1:b])
		s[b-1] = temp
	}
}
