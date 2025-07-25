/******************************************************************************/
/* decgen.go                                                                  */
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

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// encgen writes the helper functions for encoding. Intended to be
// used with go generate; see the invocation in encode.go.

// TODO: We could do more by being unsafe. Add a -unsafe flag?

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
)

var output = flag.String("output", "dec_helpers.go", "file name to write")

type Type struct {
	lower   string
	upper   string
	decoder string
}

var types = []Type{
	{
		"bool",
		"Bool",
		`slice[i] = state.decodeUint() != 0`,
	},
	{
		"complex64",
		"Complex64",
		`real := float32FromBits(state.decodeUint(), ovfl)
		imag := float32FromBits(state.decodeUint(), ovfl)
		slice[i] = complex(float32(real), float32(imag))`,
	},
	{
		"complex128",
		"Complex128",
		`real := float64FromBits(state.decodeUint())
		imag := float64FromBits(state.decodeUint())
		slice[i] = complex(real, imag)`,
	},
	{
		"float32",
		"Float32",
		`slice[i] = float32(float32FromBits(state.decodeUint(), ovfl))`,
	},
	{
		"float64",
		"Float64",
		`slice[i] = float64FromBits(state.decodeUint())`,
	},
	{
		"int",
		"Int",
		`x := state.decodeInt()
		// MinInt and MaxInt
		if x < ^int64(^uint(0)>>1) || int64(^uint(0)>>1) < x {
			error_(ovfl)
		}
		slice[i] = int(x)`,
	},
	{
		"int16",
		"Int16",
		`x := state.decodeInt()
		if x < math.MinInt16 || math.MaxInt16 < x {
			error_(ovfl)
		}
		slice[i] = int16(x)`,
	},
	{
		"int32",
		"Int32",
		`x := state.decodeInt()
		if x < math.MinInt32 || math.MaxInt32 < x {
			error_(ovfl)
		}
		slice[i] = int32(x)`,
	},
	{
		"int64",
		"Int64",
		`slice[i] = state.decodeInt()`,
	},
	{
		"int8",
		"Int8",
		`x := state.decodeInt()
		if x < math.MinInt8 || math.MaxInt8 < x {
			error_(ovfl)
		}
		slice[i] = int8(x)`,
	},
	{
		"string",
		"String",
		`u := state.decodeUint()
		n := int(u)
		if n < 0 || uint64(n) != u || n > state.b.Len() {
			errorf("length of string exceeds input size (%d bytes)", u)
		}
		if n > state.b.Len() {
			errorf("string data too long for buffer: %d", n)
		}
		// Read the data.
		data := state.b.Bytes()
		if len(data) < n {
			errorf("invalid string length %d: exceeds input size %d", n, len(data))
		}
		slice[i] = string(data[:n])
		state.b.Drop(n)`,
	},
	{
		"uint",
		"Uint",
		`x := state.decodeUint()
		/*TODO if math.MaxUint32 < x {
			error_(ovfl)
		}*/
		slice[i] = uint(x)`,
	},
	{
		"uint16",
		"Uint16",
		`x := state.decodeUint()
		if math.MaxUint16 < x {
			error_(ovfl)
		}
		slice[i] = uint16(x)`,
	},
	{
		"uint32",
		"Uint32",
		`x := state.decodeUint()
		if math.MaxUint32 < x {
			error_(ovfl)
		}
		slice[i] = uint32(x)`,
	},
	{
		"uint64",
		"Uint64",
		`slice[i] = state.decodeUint()`,
	},
	{
		"uintptr",
		"Uintptr",
		`x := state.decodeUint()
		if uint64(^uintptr(0)) < x {
			error_(ovfl)
		}
		slice[i] = uintptr(x)`,
	},
	// uint8 Handled separately.
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("decgen: ")
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("usage: decgen [--output filename]")
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by go run decgen.go -output %s; DO NOT EDIT.\n", *output)
	fmt.Fprint(&b, header)
	printMaps(&b, "Array")
	fmt.Fprint(&b, "\n")
	printMaps(&b, "Slice")
	for _, t := range types {
		fmt.Fprintf(&b, arrayHelper, t.lower, t.upper)
		fmt.Fprintf(&b, sliceHelper, t.lower, t.upper, t.decoder)
	}
	fmt.Fprintf(&b, trailer)
	source, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal("source format error:", err)
	}
	fd, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := fd.Write(source); err != nil {
		log.Fatal(err)
	}
	if err := fd.Close(); err != nil {
		log.Fatal(err)
	}
}

func printMaps(b *bytes.Buffer, upperClass string) {
	fmt.Fprintf(b, "var dec%sHelper = map[reflect.Kind]decHelper{\n", upperClass)
	for _, t := range types {
		fmt.Fprintf(b, "reflect.%s: dec%s%s,\n", t.upper, t.upper, upperClass)
	}
	fmt.Fprintf(b, "}\n")
}

const header = `
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gob

import (
	"math"
	"reflect"
)

`

const arrayHelper = `
func dec%[2]sArray(state *decoderState, v reflect.Value, length int, ovfl error) bool {
	// Can only slice if it is addressable.
	if !v.CanAddr() {
		return false
	}
	return dec%[2]sSlice(state, v.Slice(0, v.Len()), length, ovfl)
}
`

const sliceHelper = `
func dec%[2]sSlice(state *decoderState, v reflect.Value, length int, ovfl error) bool {
	slice, ok := v.Interface().([]%[1]s)
	if !ok {
		// It is kind %[1]s but not type %[1]s. TODO: We can handle this unsafely.
		return false
	}
	for i := 0; i < length; i++ {
		if state.b.Len() == 0 {
			errorf("decoding %[1]s array or slice: length exceeds input size (%%d elements)", length)
		}
		if i >= len(slice) {
			// This is a slice that we only partially allocated.
			growSlice(v, &slice, length)
		}
		%[3]s
	}
	return true
}
`

const trailer = `
// growSlice is called for a slice that we only partially allocated,
// to grow it up to length.
func growSlice[E any](v reflect.Value, ps *[]E, length int) {
	var zero E
	s := *ps
	s = append(s, zero)
	cp := cap(s)
	if cp > length {
		cp = length
	}
	s = s[:cp]
	v.Set(reflect.ValueOf(s))
	*ps = s
}
`
