//go:build editor || filedrop

/******************************************************************************/
/* window_filedrop.go                                                         */
/******************************************************************************/

package windowing

import (
	"kaijuengine.com/engine/systems/events"
	"kaijuengine.com/platform/profiler/tracing"
)

// FileDropEvent keeps a native drop as one batch so higher layers can route it
// by location before deciding how the files should be processed.
//
// NOTE: position (x, y) is relative to client area
type FileDropEvent struct {
	X     int
	Y     int
	Paths []string
}

type fileDropModule struct {
	onDrop events.EventWithArg[FileDropEvent]
}

func (w *Window) OnFileDrop() *events.EventWithArg[FileDropEvent] {
	return &w.fileDrop.onDrop
}

func (w *Window) SetFileDropEnabled(enabled bool) {
	w.setFileDropEnabled(enabled)
}

// NOTE: position (x, y) is relative to client area
func goProcessFileDropCommon(goWindow uint64, x, y int, paths []string) {
	defer tracing.NewRegion("windowing.goProcessFileDropCommon").End()
	if len(paths) == 0 {
		return
	}
	gw, ok := windowLookup.Load(goWindow)
	if !ok || gw == nil {
		return
	}
	win := gw.(*Window)
	win.fileDrop.onDrop.Execute(FileDropEvent{
		X:     x,
		Y:     y,
		Paths: paths,
	})
}
