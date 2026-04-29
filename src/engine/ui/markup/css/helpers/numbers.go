/******************************************************************************/
/* numbers.go                                                                 */
/******************************************************************************/
/*                            This file is part of                            */
/*                                KAIJU ENGINE                                */
/*                          https://kaijuengine.com/                          */
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
/* The above copyright notice and this permission notice shall be included in */
/* all copies or substantial portions of the Software.                        */
/*                                                                            */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS    */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                 */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.     */
/* IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY       */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT  */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE      */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                              */
/******************************************************************************/

package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"kaijuengine.com/rendering"
)

type WindowDimensions interface {
	DotsPerMillimeter() float64
	Width() int
	Height() int
}

var arithmeticMap = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

func ChangeNToChildCount(args []string, count int) {
	for i := range args {
		if args[i] == "n" {
			args[i] = strconv.Itoa(count)
		}
	}
}

func ArithmeticString(args []string) (int, error) {
	if len(args) == 1 {
		return strconv.Atoi(args[0])
	} else if len(args) == 2 {
		// Expected to be something like -5
		return strconv.Atoi(args[0] + args[1])
	} else {
		do := arithmeticMap["+"]
		value := 0
		negate := false
		for i := range args {
			if args[i] == "-" {
				negate = true
				continue
			} else if v, err := strconv.Atoi(args[i]); err == nil {
				if negate {
					v = -v
				}
				value = do(value, v)
			} else if f, ok := arithmeticMap[args[i]]; ok {
				do = f
			} else {
				return 0, fmt.Errorf("invalid arithmetic operator: %s", args[i])
			}
		}
		return value, nil
	}
}

func NumFromLengthWithFont(str string, window WindowDimensions, fontSize float32) float32 {
	dpmm := window.DotsPerMillimeter()
	var suffix string
	switch {
	case strings.HasSuffix(str, "vmin"):
		suffix = "vmin"
		str = str[:len(str)-4]
	case strings.HasSuffix(str, "vmax"):
		suffix = "vmax"
		str = str[:len(str)-4]
	case strings.HasSuffix(str, "rem"):
		suffix = "rem"
		str = str[:len(str)-3]
	case strings.HasSuffix(str, "vw"):
		suffix = "vw"
		str = str[:len(str)-2]
	case strings.HasSuffix(str, "vh"):
		suffix = "vh"
		str = str[:len(str)-2]
	case strings.HasSuffix(str, "ch"):
		suffix = "ch"
		str = str[:len(str)-2]
	case strings.HasSuffix(str, "px"),
		strings.HasSuffix(str, "em"),
		strings.HasSuffix(str, "ex"),
		strings.HasSuffix(str, "cm"),
		strings.HasSuffix(str, "mm"),
		strings.HasSuffix(str, "in"),
		strings.HasSuffix(str, "pt"),
		strings.HasSuffix(str, "pc"):
		suffix = str[len(str)-2:]
		str = str[:len(str)-2]
	case strings.HasSuffix(str, "%"):
		suffix = "%"
		str = str[:len(str)-1]
	}
	var size float32
	fmt.Sscanf(str, "%f", &size)
	switch suffix {
	case "%":
		size = size / 100
	case "px":
		// Read value is the size
	case "ex":
		// Relative to the font size of a lowercase letter like a, c, m, or o
		fallthrough
	case "em":
		size = size * fontSize
	case "rem":
		// TODO:
		// Root font size support is not yet wired through style inheritance.
		// For now rem is based on the engine default root em size.
		size = size * rendering.DefaultFontEMSize
	case "ch":
		// TODO:
		// Approximation until font metric support is available:
		// 1ch ~= 0.5em
		size = size * fontSize * 0.5
	case "vw":
		size = float32(window.Width()) * (size / 100)
	case "vh":
		size = float32(window.Height()) * (size / 100)
	case "vmin":
		w := float32(window.Width())
		h := float32(window.Height())
		if h < w {
			w = h
		}
		size = w * (size / 100)
	case "vmax":
		w := float32(window.Width())
		h := float32(window.Height())
		if h > w {
			w = h
		}
		size = w * (size / 100)
	case "cm":
		size = float32(dpmm) * float32(size*10)
	case "mm":
		size = float32(dpmm) * float32(size)
	case "in":
		size = float32(dpmm) * float32(size*25.4)
	case "pt":
		size = float32(dpmm) * float32(size*25.4/72)
	case "pc":
		size = float32(dpmm) * float32(size*25.4/6)
	default:
		size = 0
	}
	return size
}

// NumFromLength resolves CSS lengths with the default font size context.
// For properties that depend on the current element font, use NumFromLengthWithFont.
func NumFromLength(str string, window WindowDimensions) float32 {
	return NumFromLengthWithFont(str, window, rendering.DefaultFontEMSize)
}
