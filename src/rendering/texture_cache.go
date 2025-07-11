/******************************************************************************/
/* texture_cache.go                                                           */
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

package rendering

import (
	"kaiju/engine/assets"
	"kaiju/platform/profiler/tracing"
	"sync"
)

type TextureCache struct {
	renderer        Renderer
	assetDatabase   *assets.Database
	textures        [TextureFilterMax]map[string]*Texture
	pendingTextures []*Texture
	mutex           sync.Mutex
}

func NewTextureCache(renderer Renderer, assetDatabase *assets.Database) TextureCache {
	defer tracing.NewRegion("rendering.NewTextureCache").End()
	tc := TextureCache{
		renderer:        renderer,
		assetDatabase:   assetDatabase,
		pendingTextures: make([]*Texture, 0),
		mutex:           sync.Mutex{},
	}
	for i := range tc.textures {
		tc.textures[i] = make(map[string]*Texture)
	}
	return tc
}

func (t *TextureCache) Texture(textureKey string, filter TextureFilter) (*Texture, error) {
	defer tracing.NewRegion("TextureCache.Texture").End()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if texture, ok := t.textures[filter][textureKey]; ok {
		return texture, nil
	} else {
		if texture, err := NewTexture(t.renderer, t.assetDatabase, textureKey, filter); err == nil {
			t.pendingTextures = append(t.pendingTextures, texture)
			t.textures[filter][textureKey] = texture
			return texture, nil
		} else {
			return nil, err
		}
	}
}

func (t *TextureCache) CreatePending() {
	defer tracing.NewRegion("TextureCache.CreatePending").End()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, texture := range t.pendingTextures {
		texture.DelayedCreate(t.renderer)
	}
	t.pendingTextures = t.pendingTextures[:0]
}

func (t *TextureCache) Destroy() {
	for _, texture := range t.pendingTextures {
		texture.Destroy(t.renderer)
	}
	t.pendingTextures = t.pendingTextures[:0]
	for i := range t.textures {
		for _, texture := range t.textures[i] {
			texture.Destroy(t.renderer)
		}
		t.textures[i] = make(map[string]*Texture)
	}
}
