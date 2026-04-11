/******************************************************************************/
/* stage_workspace_hierarchy_ui.go                                            */
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

package stage_workspace

import (
	"strings"
	"weak"

	"kaijuengine.com/editor/editor_stage_manager"
	"kaijuengine.com/engine"
	"kaijuengine.com/engine/assets"
	"kaijuengine.com/engine/ui/markup/document"
	"kaijuengine.com/platform/hid"
	"kaijuengine.com/platform/profiler/tracing"
	"kaijuengine.com/platform/windowing"
	"kaijuengine.com/rendering"
)

type WorkspaceHierarchyUI struct {
	workspace            weak.Pointer[StageWorkspace]
	hierarchyArea        *document.Element
	entityTemplate       *document.Element
	entityList           *document.Element
	hierarchyDragPreview *document.Element
}

func (hui *WorkspaceHierarchyUI) setupFuncs() map[string]func(*document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.setupFuncs").End()
	return map[string]func(*document.Element){
		"hierarchySearch":        hui.hierarchySearch,
		"selectEntity":           hui.selectEntity,
		"entityToggleVisibility": hui.entityToggleVisibility,
		"entityDragStart":        hui.entityDragStart,
		"entityDrop":             hui.entityDrop,
		"entityDragEnter":        hui.entityDragEnter,
		"entityDragExit":         hui.entityDragExit,
		"hierarchyDrop":          hui.hierarchyDrop,
	}
}

func (hui *WorkspaceHierarchyUI) setup(w *StageWorkspace) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.setup").End()
	hui.hierarchyArea, _ = w.Doc.GetElementById("hierarchyArea")
	hui.entityList, _ = w.Doc.GetElementById("entityList")
	hui.entityTemplate, _ = w.Doc.GetElementById("entityTemplate")
	hui.hierarchyDragPreview, _ = w.Doc.GetElementById("hierarchyDragPreview")
	hui.workspace = weak.Make(w)
	man := w.stageView.Manager()
	man.OnEntitySpawn.Add(hui.entityCreated)
	man.OnEntityDestroy.Add(hui.entityDestroyed)
	man.OnEntitySelected.Add(hui.entitySelected)
	man.OnEntityDeselected.Add(hui.entityDeselected)
	man.OnEntityChangedParent.Add(hui.entityChangedParent)
}

func (hui *WorkspaceHierarchyUI) open() {
	defer tracing.NewRegion("WorkspaceHierarchyUI.open").End()
	hui.entityTemplate.UI.Hide()
	hui.hierarchyArea.UI.Show()
	hui.hierarchyDragPreview.UI.Hide()
}

func (hui *WorkspaceHierarchyUI) hierarchySearch(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.hierarchySearch").End()
	q := strings.ToLower(e.UI.ToInput().Text())
	for i := range hui.entityList.Children[1:] {
		lbl := hui.entityList.Children[i+1].Children[0].Children[0].UI.ToLabel()
		if strings.Contains(strings.ToLower(lbl.Text()), q) {
			hui.entityList.Children[i+1].UI.Show()
		} else {
			hui.entityList.Children[i+1].UI.Hide()
		}
	}
}

func (hui *WorkspaceHierarchyUI) processHotkeys(host *engine.Host) {
	defer tracing.NewRegion("WorkspaceContentUI.processHotkeys").End()
	kb := &host.Window.Keyboard
	if kb.KeyDown(hid.KeyboardKeyH) {
		if hui.hierarchyArea.UI.Entity().IsActive() {
			hui.hierarchyArea.UI.Hide()
		} else {
			hui.hierarchyArea.UI.Show()
		}
	} else if kb.HasCtrl() && kb.KeyDown(hid.KeyboardKeyT) {
		w := hui.workspace.Value()
		w.stageView.Manager().CreateTemplateFromSelected(w.ed.Events(), w.ed.Project())
	}
}

func (hui *WorkspaceHierarchyUI) selectEntity(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.selectEntity").End()
	id := e.Attribute("id")
	w := hui.workspace.Value()
	kb := &w.Host.Window.Keyboard
	man := w.stageView.Manager()
	if kb.HasCtrl() {
		man.SelectToggleEntityById(id)
	} else if kb.HasShift() {
		man.SelectWithChildrenOrSingleEntityById(id)
	} else {
		man.SelectEntityById(id)
	}
}

func (hui *WorkspaceHierarchyUI) textureFromString(key string) *rendering.Texture {
	w := hui.workspace.Value()
	filter := rendering.TextureFilterLinear
	tex, err := w.Host.TextureCache().Texture(key, filter)
	if err == nil {
		return tex
	}
	tex, _ = w.Host.TextureCache().Texture(assets.TextureSquare, filter)
	return tex
}

func (hui *WorkspaceHierarchyUI) entityToggleVisibility(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityToggleVisibility").End()
	id := e.Parent.Value().Attribute("id")
	w := hui.workspace.Value()
	man := w.stageView.Manager()
	if entity, ok := man.EntityById(id); ok {
		if entity.IsActive() {
			entity.Deactivate()
			w.ed.History().Add(&hierarchyEntityChangeVisibilty{
				entity:  entity,
				visible: false,
			})
		} else {
			entity.Activate()
			w.ed.History().Add(&hierarchyEntityChangeVisibilty{
				entity:  entity,
				visible: true,
			})
		}
	}
}

type HierarchyEntityDragData struct {
	hui *WorkspaceHierarchyUI
	id  string
}

func (d HierarchyEntityDragData) DragUpdate() {
	defer tracing.NewRegion("HierarchyEntityDragData.DragUpdate").End()
	m := &d.hui.workspace.Value().Host.Window.Mouse
	mp := m.ScreenPosition()
	ps := d.hui.hierarchyDragPreview.UI.Layout().PixelSize()
	d.hui.hierarchyDragPreview.UI.Layout().SetOffset(mp.X()-ps.X()*0.5, mp.Y()-ps.Y()*0.5)
}

func (hui *WorkspaceHierarchyUI) entityDragStart(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDragStart").End()
	id := e.Attribute("id")
	if id == "" {
		return
	}
	windowing.SetDragData(HierarchyEntityDragData{hui, id})
	windowing.OnDragStop.Add(hui.dragStopped)
	hui.hierarchyDragPreview.UI.Show()
}

func (hui *WorkspaceHierarchyUI) entityDrop(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDrop").End()
	dd, ok := windowing.DragData().(HierarchyEntityDragData)
	if !ok {
		return
	}
	windowing.SetDragData(nil)
	id := e.Attribute("id")
	if dd.id == id {
		return
	}
	w := hui.workspace.Value()
	man := w.stageView.Manager()
	child, ok := man.EntityById(dd.id)
	if !ok {
		return
	}
	parent, ok := man.EntityById(id)
	if !ok {
		return
	}
	man.SetEntityParent(child, parent)
	hui.clearElementDragEnterColor(e)
}

func (hui *WorkspaceHierarchyUI) entityDragEnter(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDragEnter").End()
	dd, ok := windowing.DragData().(HierarchyEntityDragData)
	if !ok {
		return
	}
	id := e.Attribute("id")
	if dd.id == id {
		return
	}
	hui.workspace.Value().Doc.SetElementClasses(
		e, hui.buildEntityClasses(e, "hierarchyEntryDragHover")...)
}

func (hui *WorkspaceHierarchyUI) entityDragExit(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDragExit").End()
	dd, ok := windowing.DragData().(HierarchyEntityDragData)
	if !ok {
		return
	}
	if dd.id == e.Attribute("id") {
		return
	}
	hui.clearElementDragEnterColor(e)
}

func (hui *WorkspaceHierarchyUI) hierarchyDrop(*document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDragExit").End()
	dd, ok := windowing.DragData().(HierarchyEntityDragData)
	if !ok {
		return
	}
	windowing.SetDragData(nil)
	w := hui.workspace.Value()
	man := w.stageView.Manager()
	child, ok := man.EntityById(dd.id)
	if !ok || child.Parent == nil {
		return
	}
	man.SetEntityParent(child, nil)
}

func (hui *WorkspaceHierarchyUI) clearElementDragEnterColor(e *document.Element) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.clearElementDragEnterColor").End()
	w := hui.workspace.Value()
	w.Doc.SetElementClasses(e, hui.buildEntityClasses(e)...)
}

func (hui *WorkspaceHierarchyUI) entityCreated(e *editor_stage_manager.StageEntity) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityCreated").End()
	w := hui.workspace.Value()
	cpy := w.Doc.DuplicateElement(hui.entityTemplate)
	w.Doc.SetElementId(cpy, e.StageData.Description.Id)
	img := cpy.Children[0].UI.ToImage()
	img.Base().ToPanel().SetUseBlending(true)
	cpy.Children[2].InnerLabel().SetText(e.Name())
	e.OnActivate.Add(func() {
		img.SetTexture(hui.textureFromString("editor/textures/icons/eye_open.png"))
	})
	e.OnDeactivate.Add(func() {
		img.SetTexture(hui.textureFromString("editor/textures/icons/eye_closed.png"))
	})
}

func (hui *WorkspaceHierarchyUI) entityDestroyed(e *editor_stage_manager.StageEntity) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDestroyed").End()
	w := hui.workspace.Value()
	if elm, ok := w.Doc.GetElementById(e.StageData.Description.Id); ok {
		hui.workspace.Value().Doc.RemoveElement(elm)
	}
}

func (hui *WorkspaceHierarchyUI) entitySelected(e *editor_stage_manager.StageEntity) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entitySelected").End()
	w := hui.workspace.Value()
	entries := w.Doc.GetElementsByClass("hierarchyEntry")
	for _, elm := range entries {
		if elm.Attribute("id") == e.StageData.Description.Id {
			hui.workspace.Value().Doc.SetElementClasses(
				elm, hui.buildEntityClasses(elm)...)
			w.Host.RunAfterNextUIClean(func() {
				hui.entityList.UI.ToPanel().ScrollToChild(elm.UI)
			})
			break
		}
	}
}

func (hui *WorkspaceHierarchyUI) entityDeselected(e *editor_stage_manager.StageEntity) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityDeselected").End()
	entries := hui.workspace.Value().Doc.GetElementsByClass("hierarchyEntry")
	for _, elm := range entries {
		if elm.Attribute("id") == e.StageData.Description.Id {
			hui.workspace.Value().Doc.SetElementClasses(
				elm, hui.buildEntityClasses(elm)...)
			break
		}
	}
}

func (hui *WorkspaceHierarchyUI) entityChangedParent(e *editor_stage_manager.StageEntity) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.entityChangedParent").End()
	w := hui.workspace.Value()
	child, ok := w.Doc.GetElementById(e.StageData.Description.Id)
	if !ok {
		return
	}
	p := editor_stage_manager.EntityToStageEntity(e.Parent)
	var parent *document.Element
	if p != nil {
		if parent, ok = w.Doc.GetElementById(p.StageData.Description.Id); !ok {
			return
		}
	} else {
		parent = hui.entityList
	}
	w.Doc.ChangeElementParent(child, parent)
	if parent.Parent.Value() == hui.entityList {
		w.Doc.SetElementClasses(parent, "hierarchyEntry")
	}
	if parent == hui.entityList {
		w.Doc.SetElementClasses(child, "hierarchyEntry")
	} else {
		hui.setIndent(child)
	}
}

func (hui *WorkspaceHierarchyUI) setIndent(row *document.Element) {
	parent := row.Parent.Value()
	if parent == nil {
		return
	}
	parentCount := 0
	for parent != hui.entityList {
		parentCount++
		parent = parent.Parent.Value()
	}
	row.Children[1].UI.ToPanel().Base().Layout().SetPadding(float32(parentCount*10), 0, 0, 0)
}

func (hui *WorkspaceHierarchyUI) dragStopped() {
	defer tracing.NewRegion("WorkspaceHierarchyUI.dragStopped").End()
	if !hui.hierarchyDragPreview.UI.Entity().IsActive() {
		return
	}
	hui.hierarchyDragPreview.UI.Hide()
}

func (hui *WorkspaceHierarchyUI) buildEntityClasses(e *document.Element, additionalClasses ...string) []string {
	defer tracing.NewRegion("WorkspaceHierarchyUI.buildEntityClasses").End()
	classes := []string{"hierarchyEntry"}
	if hui.workspace.Value().stageView.Manager().IsSelectedById(e.Attribute("id")) {
		classes = append(classes, "hierarchyEntrySelected")
	}
	classes = append(classes, additionalClasses...)
	return classes
}

func (hui *WorkspaceHierarchyUI) updateEntityName(id, name string) {
	defer tracing.NewRegion("WorkspaceHierarchyUI.updateEntityName").End()
	if e, ok := hui.workspace.Value().Doc.GetElementById(id); ok {
		e.Children[0].InnerLabel().SetText(name)
	}
}

func (hui *WorkspaceHierarchyUI) extendHeight() {
	defer tracing.NewRegion("WorkspaceHierarchyUI.extendHeight").End()
	hui.workspace.Value().Doc.SetElementClasses(hui.hierarchyArea, "edPanelBg", "sideBarTall")
}

func (hui *WorkspaceHierarchyUI) standardHeight() {
	defer tracing.NewRegion("WorkspaceHierarchyUI.standardHeight").End()
	hui.workspace.Value().Doc.SetElementClasses(hui.hierarchyArea, "edPanelBg", "sideBarStandard")
}
