/******************************************************************************/
/* rendering_tests.go                                                         */
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

package tests

import (
	"fmt"
	"kaiju/engine"
	"kaiju/engine/assets"
	"kaiju/engine/host_container"
	"kaiju/engine/systems/console"
	"kaiju/engine/ui"
	"kaiju/engine/ui/markup"
	"kaiju/engine/ui/markup/document"
	"kaiju/klib"
	"kaiju/matrix"
	"kaiju/rendering"
	"kaiju/rendering/loaders"
	"kaiju/rendering/loaders/load_result"
	"log/slog"
	"strings"
	"unsafe"
)

const TriangleShaderDataSize = int(unsafe.Sizeof(TestBasicShaderData{}))

type TestBasicShaderData struct {
	rendering.ShaderDataBase
	Color matrix.Color
}

func (t TestBasicShaderData) Size() int {
	const size = int(unsafe.Sizeof(TestBasicShaderData{}) - rendering.ShaderBaseDataStart)
	return size
}

type BoneTransform struct {
	Transform *matrix.Transform
	Skin      matrix.Mat4
}

type TestBasicSkinnedShaderData struct {
	Bones           []BoneTransform
	jointTransforms [rendering.MaxJoints]matrix.Mat4
	rendering.ShaderDataBase
	Color     matrix.Color
	SkinIndex int32
}

func (t TestBasicSkinnedShaderData) Size() int {
	const size = int(unsafe.Sizeof(TestBasicSkinnedShaderData{}) - rendering.ShaderBaseDataStart)
	return size
}

func (t *TestBasicSkinnedShaderData) NamedDataInstanceSize(name string) int {
	if name != "SkinnedUBO" {
		return 0
	}
	return int(unsafe.Sizeof(t.jointTransforms))
}

func (t *TestBasicSkinnedShaderData) UpdateNamedData(index, capacity int, name string) bool {
	if name != "SkinnedUBO" {
		return false
	}
	cap := capacity / rendering.MaxJoints / int(unsafe.Sizeof(matrix.Mat4{}))
	if index > cap {
		t.SkinIndex = int32(index % cap)
		return false
	}
	t.SkinIndex = int32(index)
	if len(t.Bones) > 0 {
		inverseRoot := t.Model()
		inverseRoot.Inverse()
		for i := range t.Bones {
			b := &t.Bones[i]
			m := matrix.Mat4Multiply(b.Skin, b.Transform.Matrix())
			parent := b.Transform.Parent()
			for parent != nil {
				m.MultiplyAssign(parent.Matrix())
				parent = parent.Parent()
			}
			t.jointTransforms[i] = m
		}
	}
	return true
}

func (t *TestBasicSkinnedShaderData) NamedDataPointer(name string) unsafe.Pointer {
	if name != "SkinnedUBO" {
		return nil
	}
	return unsafe.Pointer(&t.jointTransforms)
}

func testDrawing(uiMan *ui.Manager) {
	host := uiMan.Host
	matKey := assets.MaterialDefinitionBasic
	material, err := host.MaterialCache().Material(matKey)
	if err != nil {
		slog.Error("failed to load the material", "material", matKey, "error", err)
		return
	}
	mesh := rendering.NewMeshQuad(host.MeshCache())
	droidTex, _ := host.TextureCache().Texture("textures/android.png", rendering.TextureFilterNearest)
	tsd := TestBasicShaderData{rendering.NewShaderDataBase(), matrix.ColorWhite()}
	material = material.CreateInstance([]*rendering.Texture{droidTex})
	host.Drawings.AddDrawing(rendering.Drawing{
		Renderer:   host.Window.Renderer,
		Material:   material,
		Mesh:       mesh,
		ShaderData: &tsd,
		Transform:  nil,
	})
}

func testTwoDrawings(uiMan *ui.Manager) {
	host := uiMan.Host
	positions := []matrix.Vec3{
		{-1, 0.0, 0.0},
		{1, 0.0, 0.0},
	}
	colors := []matrix.Color{
		{1.0, 0.0, 0.0, 1.0},
		{0.0, 1.0, 0.0, 1.0},
	}
	rots := []matrix.Float{45, -45}
	matKey := assets.MaterialDefinitionBasic
	for i := 0; i < 2; i++ {
		material, err := host.MaterialCache().Material(matKey)
		if err != nil {
			slog.Error("failed to load the material", "material", matKey, "error", err)
			return
		}
		mesh := rendering.NewMeshQuad(host.MeshCache())
		droidTex, _ := host.TextureCache().Texture("textures/android.png", rendering.TextureFilterNearest)
		tsd := TestBasicShaderData{Color: colors[i]}
		m := matrix.Mat4Identity()
		m.Rotate(matrix.Vec3{0.0, rots[i], 0.0})
		m.Translate(positions[i])
		tsd.SetModel(m)
		material = material.CreateInstance([]*rendering.Texture{droidTex})
		host.Drawings.AddDrawing(rendering.Drawing{
			Renderer:   host.Window.Renderer,
			Material:   material,
			Mesh:       mesh,
			ShaderData: &tsd,
			Transform:  nil,
		})
	}
}

func testFont(uiMan *ui.Manager) {
	host := uiMan.Host
	drawings := host.FontCache().RenderMeshes(host, "Hello, World!",
		0, float32(host.Window.Height())*0.5, 0, 64, float32(host.Window.Width()), matrix.ColorBlack(), matrix.ColorDarkBG(),
		rendering.FontJustifyCenter, rendering.FontBaselineCenter,
		matrix.Vec3One(), true, false, rendering.FontRegular, 0)
	host.Drawings.AddDrawings(drawings)
}

func testOIT(uiMan *ui.Manager) {
	host := uiMan.Host
	positions := []matrix.Vec3{
		{-0.75, 0.0, -0.75},
		{-0.5, 0.0, -0.5},
		{-0.25, 0.0, -0.25},
		{0.0, 0.0, 0.0},
	}
	colors := []matrix.Color{
		{1.0, 0.0, 1.0, 0.5},
		{1.0, 0.0, 0.0, 1.0},
		{0.0, 1.0, 0.0, 0.5},
		{0.0, 0.0, 1.0, 0.5},
	}
	matKey := assets.MaterialDefinitionBasic
	material, err := host.MaterialCache().Material(matKey)
	if err != nil {
		slog.Error("failed to load the material", "material", matKey, "error", err)
		return
	}
	mesh := rendering.NewMeshQuad(host.MeshCache())
	droidTex, _ := host.TextureCache().Texture("textures/square.png", rendering.TextureFilterNearest)
	for i := 0; i < len(positions); i++ {
		tsd := TestBasicShaderData{Color: colors[i]}
		m := matrix.Mat4Identity()
		m.Translate(positions[i])
		tsd.SetModel(m)
		material = material.CreateInstance([]*rendering.Texture{droidTex})
		drawing := rendering.Drawing{
			Renderer:   host.Window.Renderer,
			Material:   material,
			Mesh:       mesh,
			ShaderData: &tsd,
			Transform:  nil,
		}
		host.Drawings.AddDrawing(drawing)
		if colors[i].A() < 1.0 {
			transparent := drawing
			m, err := host.MaterialCache().Material(assets.MaterialDefinitionBasicTransparent)
			if err != nil {
				slog.Error("failed to load the material",
					"material", assets.MaterialDefinitionBasicTransparent, "error", err)
			} else {
				transparent.Material = m
				host.Drawings.AddDrawing(transparent)
			}
		}
		host.NewEntity().SetName(fmt.Sprintf("OIT %d", i))
	}
}

func testPanel(uiMan *ui.Manager) {
	host := uiMan.Host
	tex, _ := host.TextureCache().Texture(assets.TextureSquare, rendering.TextureFilterLinear)
	p := uiMan.Add().ToPanel()
	p.Init(tex, ui.AnchorBottomLeft, ui.ElementTypePanel)
	p.DontFitContent()
	p.Base().Layout().Scale(100, 100)
	p.Base().Layout().SetOffset(10, 10)
}

func testLabel(uiMan *ui.Manager) {
	l := uiMan.Add().ToLabel()
	l.Init("Hello, World!", ui.AnchorBottomCenter)
	l.Base().Layout().Scale(100, 50)
}

func testButton(uiMan *ui.Manager) {
	host := uiMan.Host
	tex, _ := host.TextureCache().Texture(assets.TextureSquare, rendering.TextureFilterLinear)
	btn := uiMan.Add().ToButton()
	btn.Init(tex, "Click me!", ui.AnchorCenter)
	btn.Base().Layout().Scale(100, 50)
	clickCount := 0
	btn.Base().AddEvent(ui.EventTypeClick, func() {
		clickCount++
		btn.Label().SetText(fmt.Sprintf("Clicked x%d!", clickCount))
	})
}

func testHTML(uiMan *ui.Manager) {
	host := uiMan.Host
	events := map[string]func(*document.Element){
		"playGame":     func(*document.Element) { slog.Info("Clicked playGame") },
		"showSettings": func(*document.Element) { slog.Info("Clicked showSettings") },
		"showRules":    func(*document.Element) { slog.Info("Clicked showRules") },
	}
	testHTML, _ := host.AssetDatabase().ReadText("ui/tests/test.html")
	testCSS, _ := host.AssetDatabase().ReadText("ui/tests/test.css")
	markup.DocumentFromHTMLString(uiMan, testHTML, testCSS, nil, events, nil)
}

func testHTMLBinding(uiMan *ui.Manager) {
	host := uiMan.Host
	demoData := struct {
		EntityNames []string
	}{
		EntityNames: []string{"Entity 1", "\tEntity 2", "\t\tEntity 3"},
	}
	testHTML, _ := host.AssetDatabase().ReadText("ui/tests/binding.html")
	markup.DocumentFromHTMLString(uiMan, testHTML, "", demoData, nil, nil)
}

func testLayoutSimple(uiMan *ui.Manager) {
	host := uiMan.Host
	tex, _ := host.TextureCache().Texture(assets.TextureSquare, rendering.TextureFilterLinear)
	anchors := []ui.Anchor{
		ui.AnchorBottomLeft,
		ui.AnchorBottomCenter,
		ui.AnchorBottomRight,
		ui.AnchorLeft,
		ui.AnchorRight,
		ui.AnchorCenter,
		ui.AnchorTopLeft,
		ui.AnchorTopCenter,
		ui.AnchorTopRight,
	}
	panels := make([]*ui.Panel, len(anchors))
	for i := range anchors {
		panels[i] = uiMan.Add().ToPanel()
		panels[i].Init(tex, anchors[i], ui.ElementTypePanel)
	}
	for _, p := range panels {
		p.DontFitContent()
		p.Base().Layout().Scale(100, 100)
		p.Base().Layout().SetOffset(10, 10)
	}
}

func testLayout(uiMan *ui.Manager) {
	host := uiMan.Host
	tex, _ := host.TextureCache().Texture(assets.TextureSquare, rendering.TextureFilterLinear)

	p1 := uiMan.Add().ToPanel()
	p1.Init(tex, ui.AnchorTopLeft, ui.ElementTypePanel)
	p1.Base().Entity().SetName("p1")
	//p1.Layout().Scale(300, 100)

	p2 := uiMan.Add().ToPanel()
	p2.Init(tex, ui.AnchorTopLeft, ui.ElementTypePanel)
	p2.Base().Entity().SetName("p2")
	p2.SetColor(matrix.ColorBlue())
	//p2.Layout().SetPadding(5, 5, 5, 5)
	p2.Base().Layout().SetMargin(5, 5, 5, 5)
	//p2.DontFitContent()
	//p2.Layout().Scale(64, 64)
	//p2.Layout().SetOffset(10, 10)

	p3 := uiMan.Add().ToPanel()
	p3.Init(tex, ui.AnchorTopLeft, ui.ElementTypePanel)
	p3.Base().Entity().SetName("p3")
	p3.SetColor(matrix.ColorRed())
	p3.Base().Layout().Scale(32, 32)
	p3.Base().Layout().SetOffset(10, 10)
	//p3.Layout().SetMargin(5, 5, 0, 0)

	p1.AddChild(p2.Base())
	p2.AddChild(p3.Base())
}

func drawBasicMesh(host *engine.Host, res load_result.Result) {
	sd := TestBasicShaderData{rendering.NewShaderDataBase(), matrix.ColorWhite()}
	m := res.Meshes[0]
	tex, _ := host.TextureCache().Texture(assets.TextureSquare, rendering.TextureFilterLinear)
	mesh := rendering.NewMesh(m.MeshName, m.Verts, m.Indexes)
	host.MeshCache().AddMesh(mesh)
	matKey := assets.MaterialDefinitionBasic
	mat, err := host.MaterialCache().Material(matKey)
	if err != nil {
		slog.Error("failed to load the material", "material", matKey, "error", err)
		return
	}
	mat = mat.CreateInstance([]*rendering.Texture{tex})
	host.Drawings.AddDrawing(rendering.Drawing{
		Renderer:   host.Window.Renderer,
		Material:   mat,
		Mesh:       mesh,
		ShaderData: &sd,
	})
}

func testMonkeyOBJ(uiMan *ui.Manager) {
	const monkeyObj = "meshes/monkey.obj"
	host := uiMan.Host
	host.Camera.SetPosition(matrix.Vec3Backward().Scale(3))
	monkeyData := klib.MustReturn(host.AssetDatabase().ReadText(monkeyObj))
	res := loaders.OBJ(monkeyData)
	if !res.IsValid() || len(res.Meshes) != 1 {
		slog.Error("Expected 1 mesh")
		return
	}
	drawBasicMesh(host, res)
}

func testMonkeyGLTF(uiMan *ui.Manager) {
	const monkeyGLTF = "meshes/monkey.gltf"
	host := uiMan.Host
	host.Camera.SetPosition(matrix.Vec3Backward().Scale(3))
	res := klib.MustReturn(loaders.GLTF(monkeyGLTF, host.AssetDatabase()))
	if !res.IsValid() || len(res.Meshes) != 1 {
		slog.Error("Expected 1 mesh")
		return
	}
	drawBasicMesh(host, res)
}

func testMonkeyGLB(uiMan *ui.Manager) {
	const monkeyGLTF = "meshes/monkey.glb"
	host := uiMan.Host
	host.Camera.SetPosition(matrix.Vec3Backward().Scale(3))
	res := klib.MustReturn(loaders.GLTF(monkeyGLTF, host.AssetDatabase()))
	if !res.IsValid() || len(res.Meshes) != 1 {
		slog.Error("Expected 1 mesh")
		return
	}
	drawBasicMesh(host, res)
}

func testAnimationGLTF(uiMan *ui.Manager) {
	const animationGLTF = "editor/meshes/fox/Fox.gltf"
	host := uiMan.Host
	host.Camera.SetPositionAndLookAt(matrix.Vec3{150, 25, 0}, matrix.Vec3{0, 25, 0})
	//const animationGLTF = "editor/meshes/cube_animation.gltf"
	//const animationGLTF = "editor/meshes/cube_animation_slow.gltf"
	//const animationGLTF = "editor/meshes/cube_animation_slow_2.gltf"
	//host.Camera.SetPositionAndLookAt(matrix.Vec3{0, 1.5, 5}, matrix.Vec3{0, 1.5, 0})
	res := klib.MustReturn(loaders.GLTF(animationGLTF, host.AssetDatabase()))
	m := res.Meshes[0]
	textures := make([]*rendering.Texture, 0)
	for _, path := range res.Textures() {
		tex, _ := host.TextureCache().Texture(path, rendering.TextureFilterLinear)
		textures = append(textures, tex)
	}
	if len(textures) == 0 {
		tex, _ := host.TextureCache().Texture(assets.TextureSquare, rendering.TextureFilterLinear)
		textures = append(textures, tex)
	}
	entities := make([]*engine.Entity, len(res.Nodes))
	boneTransforms := make([]BoneTransform, len(res.Joints))
	for i := range res.Nodes {
		entities[i] = engine.NewEntity(uiMan.Host.WorkGroup())
		entities[i].SetName(res.Nodes[i].Name)
		entities[i].Transform = res.Nodes[i].Transform
	}
	for i := range entities {
		if res.Nodes[i].Parent >= 0 {
			entities[i].SetParent(entities[res.Nodes[i].Parent])
		}
	}
	for i := range res.Joints {
		boneTransforms[i] = BoneTransform{
			&entities[res.Joints[i].Id].Transform,
			res.Joints[i].Skin,
		}
	}
	host.AddEntities(entities...)
	mesh := rendering.NewMesh(m.MeshName, m.Verts, m.Indexes)
	host.MeshCache().AddMesh(mesh)
	sd := TestBasicSkinnedShaderData{
		Bones:     boneTransforms,
		Color:     matrix.ColorWhite(),
		SkinIndex: 0,
	}
	sd.Setup()
	matKey := assets.MaterialDefinitionBasicSkinned
	mat, err := host.MaterialCache().Material(matKey)
	if err != nil {
		slog.Error("failed to load the material", "material", matKey, "error", err)
		return
	}
	mat = mat.CreateInstance(textures)
	host.Drawings.AddDrawing(rendering.Drawing{
		Renderer:   host.Window.Renderer,
		Material:   mat,
		Mesh:       mesh,
		ShaderData: &sd,
	})
	{
		frame := 0
		animTime := 0.0
		host.Updater.AddUpdate(func(f float64) {
			animTime += f
			if animTime >= float64(res.Animations[0].Frames[frame].Time) {
				frame++
				animTime = 0
				if frame >= len(res.Animations[0].Frames) {
					frame = 0
				}
			}
			for i := range res.Animations[0].Frames[frame].Bones {
				b := &res.Animations[0].Frames[frame].Bones[i]
				var bone *matrix.Transform
				for j := range sd.Bones {
					if sd.Bones[j].Transform.Identifier == uint8(b.NodeIndex) {
						bone = sd.Bones[j].Transform
						break
					}
				}
				if bone == nil {
					continue
				}
				switch b.PathType {
				case load_result.AnimPathTranslation:
					bone.SetPosition(matrix.Vec3FromSlice(b.Data[:]))
				case load_result.AnimPathRotation:
					bone.SetRotation(matrix.Quaternion(b.Data).ToEuler())
				case load_result.AnimPathScale:
					bone.SetScale(matrix.Vec3FromSlice(b.Data[:]))
				}
			}
		})
	}
}

func SetupConsole(host *engine.Host) {
	console.For(host).AddCommand("render.test", "Open a rendering test given it's name", func(_ *engine.Host, t string) string {
		var testFunc func(*ui.Manager) = nil
		switch strings.ToLower(t) {
		case "drawing":
			testFunc = testDrawing
		case "two drawings":
			testFunc = testTwoDrawings
		case "font":
			testFunc = testFont
		case "oit":
			testFunc = testOIT
		case "panel":
			testFunc = testPanel
		case "label":
			testFunc = testLabel
		case "button":
			testFunc = testButton
		case "html":
			testFunc = testHTML
		case "layout simple":
			testFunc = testLayoutSimple
		case "layout":
			testFunc = testLayout
		case "html binding":
			testFunc = testHTMLBinding
		case "obj":
			testFunc = testMonkeyOBJ
		case "gltf":
			testFunc = testMonkeyGLTF
		case "glb":
			testFunc = testMonkeyGLB
		case "animation":
			testFunc = testAnimationGLTF
		}
		if testFunc != nil {
			c := host_container.New("Test "+t, nil)
			go c.Run(engine.DefaultWindowWidth,
				engine.DefaultWindowHeight, -1, -1)
			<-c.PrepLock
			c.Host.Camera.SetPosition(matrix.Vec3Backward().Scale(2))
			uiMan := &ui.Manager{}
			uiMan.Init(c.Host)
			testFunc(uiMan)
		}
		return "Running test"
	})
}
