package stage_workspace

import (
	"kaijuengine.com/editor/editor_stage_manager"
	"kaijuengine.com/platform/profiler/tracing"
)

type hierarchyEntityChangeVisibilty struct {
	entity  *editor_stage_manager.StageEntity
	visible bool
}

func (h *hierarchyEntityChangeVisibilty) Redo() {
	defer tracing.NewRegion("hierarchyEntityChangeVisibilty.Redo").End()
	h.entity.SetActive(h.visible)
}

func (h *hierarchyEntityChangeVisibilty) Undo() {
	defer tracing.NewRegion("hierarchyEntityChangeVisibilty.Undo").End()
	h.entity.SetActive(!h.visible)
}

func (h *hierarchyEntityChangeVisibilty) Delete() {}
func (h *hierarchyEntityChangeVisibilty) Exit()   {}
