//go:build linux && !android

/******************************************************************************/
/* window.x11.go                                                             */
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

package windowing

/*
#cgo LDFLAGS: -lX11 -lXcursor
#cgo noescape window_main
#cgo noescape window_show
#cgo noescape window_destroy
#cgo noescape window_focus
#cgo noescape window_position
#cgo noescape window_set_position
#cgo noescape window_set_size
#cgo noescape window_width_mm
#cgo noescape window_height_mm
#cgo noescape window_cursor_standard
#cgo noescape window_cursor_ibeam
#cgo noescape window_cursor_size_all
#cgo noescape window_cursor_size_ns
#cgo noescape window_cursor_size_we

#cgo nocallback window_main
#cgo nocallback window_show
#cgo nocallback window_destroy
#cgo nocallback window_focus
#cgo nocallback window_position
#cgo nocallback window_set_position
#cgo nocallback window_set_size
#cgo nocallback window_width_mm
#cgo nocallback window_height_mm
#cgo nocallback window_cursor_standard
#cgo nocallback window_cursor_ibeam
#cgo nocallback window_cursor_size_all
#cgo nocallback window_cursor_size_ns
#cgo nocallback window_cursor_size_we

#include "windowing.h"
*/
import "C"
import (
	"kaiju/klib"
	"unsafe"

	"golang.design/x/clipboard"
)

func asEventType(msg int, e *evtMem) eventType {
	switch msg {
	case 2:
		return evtKeyDown
	case 3:
		return evtKeyUp
	case 6:
		return evtMouseMove
	case 4:
		switch e.toMouseEvent().buttonId {
		case nativeMouseButtonLeft:
			return evtLeftMouseDown
		case nativeMouseButtonMiddle:
			return evtMiddleMouseDown
		case nativeMouseButtonRight:
			return evtRightMouseDown
		case nativeMouseButtonX1:
			return evtX1MouseDown
		case nativeMouseButtonX2:
			return evtX2MouseDown
		default:
			return evtUnknown
		}
	case 5:
		switch e.toMouseEvent().buttonId {
		case nativeMouseButtonLeft:
			return evtLeftMouseUp
		case nativeMouseButtonMiddle:
			return evtMiddleMouseUp
		case nativeMouseButtonRight:
			return evtRightMouseUp
		case nativeMouseButtonX1:
			return evtX1MouseUp
		case nativeMouseButtonX2:
			return evtX2MouseUp
		default:
			return evtUnknown
		}
	case 9:
		fallthrough
	case 10:
		return evtActivity
	default:
		return evtUnknown
	}
}

func scaleScrollDelta(delta float32) float32 {
	return delta
}

func createWindow(windowName string, width, height, x, y int, evtSharedMem *evtMem) {
	title := C.CString(windowName)
	C.window_main(title, C.int(width), C.int(height),
		C.int(x), C.int(y), evtSharedMem.AsPointer(), evtSharedMemSize)
	C.free(unsafe.Pointer(title))
}

func (w *Window) showWindow(evtSharedMem *evtMem) {
	C.window_show(w.handle)
}

func (w *Window) destroy() {
	C.window_destroy(w.handle)
}

func (w *Window) poll() {
	//evtType := uint32(C.window_poll_controller(w.handle))
	//if evtType != 0 {
	//	w.processControllerEvent(asEventType(evtType))
	//}
	evtType := 1
	for evtType != 0 && !w.evtSharedMem.IsQuit() {
		evtType = int(C.window_poll(w.handle))
		if evtType != 0 {
			t := asEventType(evtType, w.evtSharedMem)
			w.processEvent(t)
		}
	}
}

func (w *Window) cursorStandard() {
	C.window_cursor_standard(w.handle)
}

func (w *Window) cursorIbeam() {
	C.window_cursor_ibeam(w.handle)
}

func (w *Window) cursorSizeAll() {
	C.window_cursor_size_all(w.handle)
}

func (w *Window) cursorSizeNS() {
	C.window_cursor_size_ns(w.handle)
}

func (w *Window) cursorSizeWE() {
	C.window_cursor_size_we(w.handle)
}

func (w *Window) copyToClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}

func (w *Window) clipboardContents() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func (w *Window) sizeMM() (int, int, error) {
	width := C.window_width_mm(w.handle)
	height := C.window_height_mm(w.handle)
	return int(width), int(height), nil
}

func (w *Window) cHandle() unsafe.Pointer   { return C.window(w.handle) }
func (w *Window) cInstance() unsafe.Pointer { return C.display(w.handle) }

func (w *Window) focus() {
	C.window_focus(w.handle)
}

func (w *Window) position() (x, y int) {
	C.window_position(w.handle, (*C.int)(unsafe.Pointer(&x)), (*C.int)(unsafe.Pointer(&y)))
	return x, y
}

func (w *Window) setPosition(x, y int) {
	C.window_set_position(w.handle, C.int(x), C.int(y))
}

func (w *Window) setSize(width, height int) {
	C.window_set_size(w.handle, C.int(width), C.int(height))
}

func (w *Window) removeBorder() {
	klib.NotYetImplemented(234)
}

func (w *Window) addBorder() {
	klib.NotYetImplemented(234)
}
