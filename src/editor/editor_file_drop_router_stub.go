//go:build !editor && !filedrop

/******************************************************************************/
/* editor_file_drop_router_stub.go                                            */
/******************************************************************************/

package editor

type FileDropRouter struct{}

func (ed *Editor) FileDropRouter() *FileDropRouter { return &ed.fileDropRouter }

func (ed *Editor) connectFileDropRouter() {}
