//go:build !editor && !filedrop

/******************************************************************************/
/* window_filedrop_stub.go                                                    */
/******************************************************************************/

package windowing

type fileDropModule struct{}

func (m *fileDropModule) processQueuedFileDrops() {}
