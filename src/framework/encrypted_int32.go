/******************************************************************************/
/* encrypted_int32.go                                                         */
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

package framework

import "math/rand/v2"

type EncryptedInt32 struct {
	RawValue int32
	seed     int32
}

func (e EncryptedInt32) Value() int32 {
	return e.RawValue ^ e.seed
}

func (e *EncryptedInt32) SetValue(value int32) {
	e.seed = rand.Int32()
	e.RawValue = value ^ e.seed
}

func (e *EncryptedInt32) Increment() {
	e.SetValue(e.Value() + 1)
}

func (e *EncryptedInt32) Decrement() {
	e.SetValue(e.Value() - 1)
}

func (e *EncryptedInt32) Add(amount int32) {
	e.SetValue(e.Value() + amount)
}

func (e *EncryptedInt32) Subtract(amount int32) {
	e.SetValue(e.Value() - amount)
}

func (e *EncryptedInt32) Decrypt() {
	e.RawValue ^= e.seed
	e.seed = 0
}

func (e *EncryptedInt32) Encrypt() {
	e.SetValue(e.RawValue)
}
