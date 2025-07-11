/******************************************************************************/
/* host.go                                                                    */
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

package engine

import (
	"kaiju/engine/assets"
	"kaiju/engine/cameras"
	"kaiju/engine/collision_system"
	"kaiju/engine/systems/events"
	"kaiju/engine/systems/logging"
	"kaiju/engine/systems/tweening"
	"kaiju/klib"
	"kaiju/matrix"
	"kaiju/platform/audio"
	"kaiju/platform/concurrent"
	"kaiju/platform/profiler/tracing"
	"kaiju/platform/windowing"
	"kaiju/plugins"
	"kaiju/rendering"
	"log/slog"
	"math"
	"slices"
	"time"
)

// FrameId is a unique identifier for a frame
type FrameId = uint64

// InvalidFrameId can be used to indicate that a frame id is invalid
const InvalidFrameId = math.MaxUint64

type frameRun struct {
	frame FrameId
	call  func()
}

type timeRun struct {
	end  time.Time
	call func()
}

// Host is the mediator to the entire runtime for the game/editor. It is the
// main entry point for the game loop and is responsible for managing all
// entities, the window, and the rendering context. The host can be used to
// create and manage entities, call update functions on the main thread, and
// access various caches and resources.
//
// The host is expected to be passed around quite often throughout the program.
// It is designed to remove things like service locators, singletons, and other
// global state. You can have multiple hosts in a program to isolate things like
// windows and game state.
type Host struct {
	name             string
	editorEntities   editorEntities
	entities         []*Entity
	entityLookup     map[EntityId]*Entity
	lights           []rendering.Light
	timeRunner       []timeRun
	frameRunner      []frameRun
	plugins          []*plugins.LuaVM
	Window           *windowing.Window
	LogStream        *logging.LogStream
	workGroup        concurrent.WorkGroup
	threads          concurrent.Threads
	Camera           cameras.Camera
	UICamera         cameras.Camera
	collisionManager collision_system.Manager
	audio            *audio.Audio
	shaderCache      rendering.ShaderCache
	textureCache     rendering.TextureCache
	meshCache        rendering.MeshCache
	fontCache        rendering.FontCache
	materialCache    rendering.MaterialCache
	Drawings         rendering.Drawings
	frame            FrameId
	frameTime        float64
	Closing          bool
	UIUpdater        Updater
	UILateUpdater    Updater
	Updater          Updater
	LateUpdater      Updater
	assetDatabase    assets.Database
	OnClose          events.Event
	CloseSignal      chan struct{}
	frameRateLimit   *time.Ticker
	inEditorEntity   int
}

// NewHost creates a new host with the given name and log stream. The log stream
// is the log handler that is used by the slog package functions. A Host that
// is created through NewHost has no function until #Host.Initialize is called.
//
// This is primarily called from #host_container/New
func NewHost(name string, logStream *logging.LogStream) *Host {
	w := float32(DefaultWindowWidth)
	h := float32(DefaultWindowHeight)
	host := &Host{
		name:          name,
		entities:      make([]*Entity, 0),
		frameTime:     0,
		Closing:       false,
		UIUpdater:     NewUpdater(),
		UILateUpdater: NewUpdater(),
		Updater:       NewUpdater(),
		LateUpdater:   NewUpdater(),
		assetDatabase: assets.NewDatabase(),
		Drawings:      rendering.NewDrawings(),
		CloseSignal:   make(chan struct{}, 1),
		Camera:        cameras.NewStandardCamera(w, h, w, h, matrix.Vec3Backward()),
		UICamera:      cameras.NewStandardCameraOrthographic(w, h, w, h, matrix.Vec3{0, 0, 250}),
		LogStream:     logStream,
		entityLookup:  make(map[EntityId]*Entity),
		threads:       concurrent.NewThreads(),
	}
	return host
}

// Initializes the various systems and caches that are mediated through the
// host. This includes the window, the shader cache, the texture cache, the mesh
// cache, and the font cache, and the camera systems.
func (host *Host) Initialize(width, height, x, y int) error {
	if width <= 0 {
		width = DefaultWindowWidth
	}
	if height <= 0 {
		height = DefaultWindowHeight
	}
	win, err := windowing.New(host.name, width, height, x, y, &host.assetDatabase)
	if err != nil {
		return err
	}
	host.Window = win
	host.threads.Start()
	host.Camera.ViewportChanged(float32(width), float32(height))
	host.UICamera.ViewportChanged(float32(width), float32(height))
	host.shaderCache = rendering.NewShaderCache(host.Window.Renderer, &host.assetDatabase)
	host.textureCache = rendering.NewTextureCache(host.Window.Renderer, &host.assetDatabase)
	host.meshCache = rendering.NewMeshCache(host.Window.Renderer, &host.assetDatabase)
	host.fontCache = rendering.NewFontCache(host.Window.Renderer, &host.assetDatabase)
	host.materialCache = rendering.NewMaterialCache(host.Window.Renderer, &host.assetDatabase)
	host.Window.OnResize.Add(host.resized)
	return nil
}

func (host *Host) InitializeRenderer() error {
	w, h := int32(host.Window.Width()), int32(host.Window.Height())
	if err := host.Window.Renderer.Initialize(host, w, h); err != nil {
		slog.Error("failed to initialize the renderer", "error", err)
		return err
	}
	if err := host.FontCache().Init(host.Window.Renderer, host.AssetDatabase(), host); err != nil {
		slog.Error("failed to initialize the font cache", "error", err)
		return err
	}
	if err := rendering.SetupLightMaterials(host.MaterialCache()); err != nil {
		slog.Error("failed to setup the light materials", "error", err)
		return err
	}
	return nil
}

func (host *Host) InitializeAudio() (err error) {
	host.audio, err = audio.New()
	return err
}

// WorkGroup returns the work group for this instance of host
func (host *Host) WorkGroup() *concurrent.WorkGroup { return &host.workGroup }

// Threads returns the long-running threads for this instance of host
func (host *Host) Threads() *concurrent.Threads { return &host.threads }

// Name returns the name of the host
func (host *Host) Name() string { return host.name }

// CreatingEditorEntities is used exclusively for the editor to know that the
// entities that are being created are for the editor. This is used to logically
// separate editor entities from game entities.
//
// This will increment so it can be called many times, however it is expected
// that #Host.DoneCreatingEditorEntities is be called the same number of times.
func (host *Host) CreatingEditorEntities() {
	host.inEditorEntity++
}

// DoneCreatingEditorEntities is used to signal that the editor is done creating
// entities. This should be called the same number of times as
// #Host.CreatingEditorEntities. When the internal counter reaches 0, then any
// entity created on the host will go to the standard entity pool.
func (host *Host) DoneCreatingEditorEntities() {
	host.inEditorEntity--
}

// CollisionManager returns the collision manager for this host
func (host *Host) CollisionManager() *collision_system.Manager {
	return &host.collisionManager
}

// ShaderCache returns the shader cache for the host
func (host *Host) ShaderCache() *rendering.ShaderCache {
	return &host.shaderCache
}

// TextureCache returns the texture cache for the host
func (host *Host) TextureCache() *rendering.TextureCache {
	return &host.textureCache
}

// MeshCache returns the mesh cache for the host
func (host *Host) MeshCache() *rendering.MeshCache {
	return &host.meshCache
}

// FontCache returns the font cache for the host
func (host *Host) FontCache() *rendering.FontCache {
	return &host.fontCache
}

// MaterialCache returns the font cache for the host
func (host *Host) MaterialCache() *rendering.MaterialCache {
	return &host.materialCache
}

// AssetDatabase returns the asset database for the host
func (host *Host) AssetDatabase() *assets.Database {
	return &host.assetDatabase
}

// Plugins returns all of the loaded plugins for the host
func (host *Host) Plugins() []*plugins.LuaVM {
	return host.plugins
}

// Audio returns the audio system for the host
func (host *Host) Audio() *audio.Audio {
	return host.audio
}

// ClearEntities will remove all entities from the host. This will remove all
// entities from the standard entity pool only. The entities will be destroyed
// using the standard destroy method, so they will take not be fully removed
// during the frame that this function was called.
func (host *Host) ClearEntities() {
	for _, e := range host.entities {
		e.Destroy()
	}
}

// RemoveEntity removes an entity from the host. This will remove the entity
// from the standard entity pool. This will determine if the entity is in the
// editor entity pool and remove it from there if so, otherwise it will be
// removed from the standard entity pool. Entities are not ordered, so they are
// removed in O(n) time. Do not assume the entities are ordered at any time.
func (host *Host) RemoveEntity(entity *Entity) {
	if host.editorEntities.contains(entity) {
		host.editorEntities.remove(entity)
	} else {
		for i, e := range host.entities {
			if e == entity {
				host.entities = klib.RemoveUnordered(host.entities, i)
				break
			}
		}
	}
}

// AddEntity adds an entity to the host. This will add the entity to the
// standard entity pool. If the host is in the process of creating editor
// entities, then the entity will be added to the editor entity pool.
func (host *Host) AddEntity(entity *Entity) {
	host.addEntity(entity)
}

// AddEntities adds multiple entities to the host. This will add the entities
// using the same rules as AddEntity. If the host is in the process of creating
// editor entities, then the entities will be added to the editor entity pool.
func (host *Host) AddEntities(entities ...*Entity) {
	host.addEntities(entities...)
}

// AddLight adds a light to the internal list of lights the host is aware of
func (host *Host) AddLight(light rendering.Light) {
	host.lights = append(host.lights, light)
}

// Lights returns all of the active lights managed by this host
func (host *Host) Lights() []rendering.Light {
	return host.lights
}

// ClearLights clears out all of the lights that the host is tracking
func (host *Host) ClearLights() {
	host.lights = host.lights[:0]
}

// FindEntity will search for an entity contained in this host by its id. If the
// entity is found, then it will return the entity and true, otherwise it will
// return nil and false.
func (host *Host) FindEntity(id EntityId) (*Entity, bool) {
	e, ok := host.entityLookup[id]
	return e, ok
}

// Entities returns all the entities that are currently in the host. This will^
// return all entities in the standard entity pool only. In the editor, this
// will not return any entities that have been destroyed (and are pending
// cleanup due to being in the undo history)
func (host *Host) Entities() []*Entity { return host.selectAllValidEntities() }

// Entities returns all the entities that are currently in the host. This will
// return all entities in the standard entity pool only. In the editor, this
// will also return any entities that have been destroyed (and are pending
// cleanup due to being in the undo history)
func (host *Host) EntitiesRaw() []*Entity { return host.entities }

// NewEntity creates a new entity and adds it to the host. This will add the
// entity to the standard entity pool. If the host is in the process of creating
// editor entities, then the entity will be added to the editor entity pool.
func (host *Host) NewEntity() *Entity {
	entity := NewEntity(&host.workGroup)
	host.AddEntity(entity)
	return entity
}

// Update is the main update loop for the host. This will poll the window for
// events, update the entities, and render the scene. This will also check if
// the window has been closed or crashed and set the closing flag accordingly.
//
// The update order is FrameRunner -> Update -> LateUpdate -> EndUpdate:
//
// [-] FrameRunner: Functions added to RunAfterFrames
// [-] UIUpdate: Functions added to UIUpdater
// [-] UILateUpdate: Functions added to UILateUpdater
// [-] Update: Functions added to Updater
// [-] LateUpdate: Functions added to LateUpdater
// [-] EndUpdate: Internal functions for preparing for the next frame
//
// Any destroyed entities will also be ticked for their cleanup. This will also
// tick the editor entities for cleanup.
func (host *Host) Update(deltaTime float64) {
	defer tracing.NewRegion("Host.Update").End()
	host.frame++
	host.frameTime += deltaTime
	host.Window.Poll()
	for i := 0; i < len(host.frameRunner); i++ {
		if host.frameRunner[i].frame <= host.frame {
			host.frameRunner[i].call()
			host.frameRunner = klib.RemoveUnordered(host.frameRunner, i)
			i--
		}
	}
	if len(host.timeRunner) > 0 {
		now := time.Now()
		for i := 0; i < len(host.timeRunner); i++ {
			if host.timeRunner[i].end.Before(now) {
				host.timeRunner[i].call()
				host.timeRunner = klib.RemoveUnordered(host.timeRunner, i)
				i--
			}
		}
	}
	host.UIUpdater.Update(deltaTime)
	host.UILateUpdater.Update(deltaTime)
	tweening.Update(deltaTime)
	host.Updater.Update(deltaTime)
	host.LateUpdater.Update(deltaTime)
	host.collisionManager.Update(deltaTime)
	if host.Window.IsClosed() || host.Window.IsCrashed() {
		host.Closing = true
	}
	end := len(host.entities)
	back := end
	for i := 0; i < back; i++ {
		e := host.entities[i]
		if e.TickCleanup() {
			host.entities[i] = host.entities[back-1]
			back--
			i--
		}
	}
	host.entities = host.entities[:back]
	host.editorEntities.tickCleanup()
	host.Window.EndUpdate()
}

// Render will render the scene. This starts by preparing any drawings that are
// pending. It also creates any pending shaders, textures, and meshes before
// the start of the render. The frame is then readied, buffers swapped, and any
// transformations that are dirty on entities are then cleaned.
func (host *Host) Render() {
	defer tracing.NewRegion("Host.Render").End()
	host.workGroup.Execute(matrix.TransformWorkGroup, &host.threads)
	host.Drawings.PreparePending()
	host.shaderCache.CreatePending()
	host.textureCache.CreatePending()
	host.meshCache.CreatePending()
	if host.Drawings.HasDrawings() {
		if host.Window.Renderer.ReadyFrame(host.Camera,
			host.UICamera, host.lights, float32(host.Runtime())) {
			host.Drawings.Render(host.Window.Renderer)
		}
	}
	host.Window.SwapBuffers()
	host.workGroup.Execute(matrix.TransformResetWorkGroup, &host.threads)
	//host.editorEntities.resetDirty()
}

// Frame will return the current frame id
func (host *Host) Frame() FrameId { return host.frame }

// Runtime will return how long the host has been running in seconds
func (host *Host) Runtime() float64 { return host.frameTime }

// RunAfterFrames will call the given function after the given number of frames
// have passed from the current frame
func (host *Host) RunAfterFrames(wait int, call func()) {
	host.frameRunner = append(host.frameRunner, frameRun{
		frame: host.frame + uint64(wait),
		call:  call,
	})
}

func (host *Host) RunOnMainThread(call func()) {
	host.frameRunner = append(host.frameRunner, frameRun{
		frame: host.frame,
		call:  call,
	})
}

// RunAfterTime will call the given function after the given number of time
// has passed from the current frame
func (host *Host) RunAfterTime(wait time.Duration, call func()) {
	host.timeRunner = append(host.timeRunner, timeRun{
		end:  time.Now().Add(wait),
		call: call,
	})
}

// Teardown will destroy the host and all of its resources. This will also
// execute the OnClose event. This will also signal the CloseSignal channel.
func (host *Host) Teardown() {
	host.Window.Renderer.WaitForRender()
	host.OnClose.Execute()
	host.UIUpdater.Destroy()
	host.UILateUpdater.Destroy()
	host.Updater.Destroy()
	host.LateUpdater.Destroy()
	host.Drawings.Destroy(host.Window.Renderer)
	host.textureCache.Destroy()
	host.meshCache.Destroy()
	host.shaderCache.Destroy()
	host.fontCache.Destroy()
	host.materialCache.Destroy()
	host.assetDatabase.Destroy()
	host.Window.Destroy()
	host.threads.Stop()
	host.CloseSignal <- struct{}{}
}

// WaitForFrameRate will block until the desired frame rate limit is reached
func (h *Host) WaitForFrameRate() {
	defer tracing.NewRegion("Host.WaitForFrameRate").End()
	if h.frameRateLimit != nil {
		<-h.frameRateLimit.C
	}
}

// SetFrameRateLimit will set the frame rate limit for the host. If the frame
// rate is set to 0, then the frame rate limit will be removed.
//
// If a frame rate is set, then the host will block until the desired frame rate
// is reached before continuing the update loop.
func (h *Host) SetFrameRateLimit(fps int64) {
	defer tracing.NewRegion("Host.SetFrameRateLimit").End()
	if fps == 0 {
		h.frameRateLimit.Stop()
		h.frameRateLimit = nil
	} else {
		h.frameRateLimit = time.NewTicker(time.Second / time.Duration(fps))
	}
}

// Close will set the closing flag to true and signal the host to clean up
// resources and close the window.
func (host *Host) Close() {
	host.Closing = true
}

// ReserveEntities will grow the internal entity list by the given amount,
// this is useful for when you need to create a large amount of entities
func (host *Host) ReserveEntities(additional int) {
	defer tracing.NewRegion("Host.ReserveEntities").End()
	host.entities = slices.Grow(host.entities, additional)
}

// BootstrapPlugins will initialize the plugin interface and read all of the
// plugins that are in the content folder and prepare them for execution
func (host *Host) BootstrapPlugins() error {
	defer tracing.NewRegion("Host.BootstrapPlugins").End()
	var err error
	host.plugins, err = plugins.LaunchPlugins(host.AssetDatabase())
	return err
}

func (host *Host) resized() {
	w, h := float32(host.Window.Width()), float32(host.Window.Height())
	host.Camera.ViewportChanged(w, h)
	host.UICamera.ViewportChanged(w, h)
}
