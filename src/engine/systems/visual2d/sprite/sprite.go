/******************************************************************************/
/* sprite.go                                                                  */
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

package sprite

import (
	"kaiju/engine"
	"kaiju/engine/assets"
	"kaiju/matrix"
	"kaiju/platform/profiler/tracing"
	"kaiju/rendering"
	"log/slog"
)

var ZAxisScaleFactor = float32(16.0)

type Sprite struct {
	Entity                   *engine.Entity
	host                     *engine.Host
	texture                  *rendering.Texture
	flipBook                 []*rendering.Texture
	frameDelay, fps          float32
	frameCount, currentFrame int
	paused                   bool
	spriteSheet              SpriteSheet
	shaderData               ShaderData
	drawing                  rendering.Drawing
	currentClipName          string
	currentClip              *SpriteSheetClip
	clipIdx                  int
	baseScale                matrix.Vec3
}

func (s Sprite) isFlipBook() bool {
	return len(s.flipBook) > 0
}

func (s Sprite) isSpriteSheet() bool {
	return len(s.spriteSheet.clipList) > 0
}

func (s *Sprite) Resize(width, height matrix.Float) {
	s.Entity.Transform.SetScale(matrix.Vec3{width, height, 1.0})
}

func (s *Sprite) Position() matrix.Vec2 {
	return s.Entity.Transform.Position().AsVec2()
}

func (s *Sprite) Move(x, y matrix.Float) {
	p := s.Entity.Transform.Position()
	s.Entity.Transform.SetPosition(matrix.Vec3{p.X() + x, p.Y() + y, p.Z()})
}

func (s *Sprite) Move3D(x, y, z matrix.Float) {
	p := s.Entity.Transform.Position()
	s.Entity.Transform.SetPosition(matrix.Vec3{p.X() + x, p.Y() + y, p.Z() + z})
}

func (s *Sprite) SetPosition(x, y matrix.Float) {
	p := s.Entity.Transform.Position()
	x -= matrix.Float(s.host.Window.Width()) * 0.5
	y -= matrix.Float(s.host.Window.Height()) * 0.5
	s.Entity.Transform.SetPosition(matrix.Vec3{x, y, p.Z()})
}

func (s *Sprite) SetPosition3D(x, y, z matrix.Float) {
	s.Entity.Transform.SetPosition(matrix.Vec3{x, y, z})
	scale := s.baseScale.Add(matrix.Vec3One().Scale(z * ZAxisScaleFactor))
	scale.SetZ(1)
	s.Entity.Transform.SetScale(scale)
}

func (s *Sprite) PlayAnimation() {
	s.paused = false
}

func (s *Sprite) StopAnimation() {
	s.paused = true
}

func (s *Sprite) SetFrameRate(framesPerSecond float32) {
	s.fps = framesPerSecond
}

func (s *Sprite) recreateDrawing() {
	defer tracing.NewRegion("Sprite.recreateDrawing").End()
	s.shaderData.Destroy()
	proxy := s.shaderData
	proxy.CancelDestroy()
	s.shaderData = proxy
}

func (s *Sprite) SetTexture(texture *rendering.Texture) {
	defer tracing.NewRegion("Sprite.SetTexture").End()
	s.texture = texture
	if s.drawing.IsValid() {
		s.recreateDrawing()
		s.drawing.Material.Textures[0] = texture
		s.host.Drawings.AddDrawing(s.drawing)
	}
}

func (s *Sprite) SetFlipBookAnimation(framesPerSecond float32, textures ...*rendering.Texture) {
	defer tracing.NewRegion("Sprite.SetFlipBookAnimation").End()
	count := len(textures)
	s.flipBook = make([]*rendering.Texture, 0, count)
	for i := 0; i < count; i++ {
		s.flipBook = append(s.flipBook, textures[i])
	}
	s.frameCount = count
	s.fps = framesPerSecond
	s.currentFrame = 0
	s.resetDelay()
	s.SetTexture(s.flipBook[s.currentFrame])
}

func (s *Sprite) SetColor(color matrix.Color) {
	defer tracing.NewRegion("Sprite.SetColor").End()
	s.shaderData.FgColor = color
	if color.A() < 1 {
		s.recreateDrawing()
		s.host.Drawings.AddDrawing(s.drawing)
	}
}

func (s Sprite) CurrentClipName() string { return s.currentClipName }

func (s *Sprite) SetSheetClip(clipName string) {
	defer tracing.NewRegion("Sprite.SetSheetClip").End()
	if s.currentClipName != clipName {
		s.currentClipName = clipName
		s.currentClip = s.spriteSheet.clips[clipName]
		s.setSheetFrame(0)
		s.frameCount = len(s.currentClip.Frames)
	}
}

func (s *Sprite) resetDelay() {
	s.frameDelay = 1.0 / s.fps
}

func (s *Sprite) update(deltaTime float64) {
	defer tracing.NewRegion("Sprite.update").End()
	if !s.Entity.IsActive() {
		return
	}
	s.frameDelay -= float32(deltaTime)
	if s.frameCount > 0 && s.frameDelay <= 0.0 {
		s.currentFrame++
		if s.currentFrame == s.frameCount {
			s.currentFrame = 0
		}
		if s.isFlipBook() {
			s.SetTexture(s.flipBook[s.currentFrame])
		} else if s.isSpriteSheet() {
			frame := s.clipIdx + 1
			if frame < s.frameCount {
				s.setSheetFrame(frame)
			} else {
				s.setSheetFrame(0)
			}
		}
		s.resetDelay()
	}
}

func (s *Sprite) setSheetFrame(frame int) {
	defer tracing.NewRegion("Sprite.setSheetFrame").End()
	s.clipIdx = frame
	f := s.currentClip.Frames[frame]
	h := float32(f.Height()) / s.texture.Size().Height()
	s.shaderData.UVs = matrix.Vec4{
		float32(f.Left()) / s.texture.Size().Width(),
		1.0 - h - float32(f.Top())/s.texture.Size().Height(),
		float32(f.Width()) / s.texture.Size().Width(),
		h,
	}
}

func (s *Sprite) Activate() {
	defer tracing.NewRegion("Sprite.Activate").End()
	s.Entity.Activate()
	s.shaderData.Activate()
}

func (s *Sprite) Deactivate() {
	defer tracing.NewRegion("Sprite.Deactivate").End()
	s.Entity.Deactivate()
	s.shaderData.Deactivate()
}

func NewSprite(x, y, width, height matrix.Float, host *engine.Host, texture *rendering.Texture, color matrix.Color) *Sprite {
	defer tracing.NewRegion("sprite.NewSprite").End()
	e := host.NewEntity()
	sprite := &Sprite{
		host:     host,
		Entity:   e,
		flipBook: []*rendering.Texture{},
	}
	sprite.baseScale = matrix.Vec3{width, height, 1.0}
	mat, err := host.MaterialCache().Material(assets.MaterialDefinitionSprite)
	if err != nil {
		slog.Error("failed to load the sprite material", "error", err)
		return nil
	}
	mesh := rendering.NewMeshQuad(host.MeshCache())
	sprite.Entity.Transform.SetPosition(matrix.Vec3{x, y, 0})
	sprite.Entity.Transform.SetScale(matrix.Vec3{width, height, 1})
	sprite.texture = texture
	sprite.shaderData = ShaderData{
		ShaderDataBase: rendering.NewShaderDataBase(),
		FgColor:        color,
		UVs:            matrix.Vec4{0.0, 0.0, 1.0, 1.0},
	}
	drawing := rendering.Drawing{
		Renderer:   host.Window.Renderer,
		Material:   mat.CreateInstance([]*rendering.Texture{texture}),
		Mesh:       mesh,
		ShaderData: &sprite.shaderData,
		Transform:  &sprite.Entity.Transform,
	}
	host.Drawings.AddDrawing(drawing)
	if color.A() < 1 {
		transparent := drawing
		m, err := host.MaterialCache().Material(assets.MaterialDefinitionSpriteTransparent)
		if err != nil {
			slog.Error("failed to load the material",
				"material", assets.MaterialDefinitionSpriteTransparent, "error", err)
		} else {
			transparent.Material = m.CreateInstance(drawing.Material.Textures)
			host.Drawings.AddDrawing(transparent)
		}
	}
	return sprite
}

func NewSpriteFlipBook(x, y, width, height float32, host *engine.Host, images []*rendering.Texture, fps float32) *Sprite {
	defer tracing.NewRegion("sprite.NewSpriteFlipBook").End()
	s := NewSprite(x, y, width, height, host, images[0], matrix.ColorWhite())
	s.fps = fps
	s.flipBook = images
	updateId := host.Updater.AddUpdate(s.update)
	s.Entity.OnDestroy.Add(func() {
		host.Updater.RemoveUpdate(updateId)
	})
	return s
}

func NewSpriteSheet(x, y, width, height float32, host *engine.Host, texture *rendering.Texture, jsonStr string, fps float32, initialClip string) *Sprite {
	defer tracing.NewRegion("sprite.NewSpriteSheet").End()
	s := NewSprite(x, y, width, height, host, texture, matrix.ColorWhite())
	var err error
	s.spriteSheet, err = ReadSpriteSheetData(jsonStr)
	if err != nil {
		panic(err)
	}
	s.fps = fps
	s.SetSheetClip(initialClip)
	updateId := host.Updater.AddUpdate(s.update)
	s.Entity.OnDestroy.Add(func() {
		host.Updater.RemoveUpdate(updateId)
	})
	return s
}
