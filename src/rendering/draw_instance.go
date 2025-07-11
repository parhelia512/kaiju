/******************************************************************************/
/* draw_instance.go                                                           */
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
	"kaiju/engine/runtime/encoding/gob"
	"kaiju/klib"
	"kaiju/matrix"
	"reflect"
	"unsafe"
)

func init() {
	gob.Register(&ShaderDataBasic{})
}

type DrawInstance interface {
	Destroy()
	IsDestroyed() bool
	Activate()
	Deactivate()
	IsActive() bool
	Size() int
	SetModel(model matrix.Mat4)
	UpdateModel()
	DataPointer() unsafe.Pointer
	// Returns true if it should write the data, otherwise false
	UpdateNamedData(index, capacity int, name string) bool
	NamedDataPointer(name string) unsafe.Pointer
	NamedDataInstanceSize(name string) int
	setTransform(transform *matrix.Transform)
	setShadow(shadow DrawInstance)
}

func ReflectDuplicateDrawInstance(target DrawInstance) DrawInstance {
	val := reflect.ValueOf(target)
	if !val.IsValid() {
		return nil
	}
	newVal := reflect.New(val.Elem().Type()).Elem()
	newVal.Set(val.Elem())
	return newVal.Addr().Interface().(DrawInstance)
}

const ShaderBaseDataStart = unsafe.Offsetof(ShaderDataBase{}.model)

type ShaderDataBase struct {
	destroyed   bool
	deactivated bool
	_           [2]byte
	shadow      DrawInstance
	transform   *matrix.Transform
	InitModel   matrix.Mat4
	model       matrix.Mat4
}

type ShaderDataUnlit struct {
	ShaderDataBase
	Color matrix.Color
	UVs   matrix.Vec4
}

func (t ShaderDataUnlit) Size() int {
	return int(unsafe.Sizeof(ShaderDataUnlit{}) - ShaderBaseDataStart)
}

type ShaderDataBasic struct {
	ShaderDataBase
	Color matrix.Color
}

func (t ShaderDataBasic) Size() int {
	return int(unsafe.Sizeof(ShaderDataBasic{}) - ShaderBaseDataStart)
}

type ShaderDataPBR struct {
	ShaderDataBase
	VertColors matrix.Color
	Metallic   float32
	Roughness  float32
	Emissive   float32
	Light0     float32
	Light1     float32
	Light2     float32
	Light3     float32
}

func (t ShaderDataPBR) Size() int {
	return int(unsafe.Sizeof(ShaderDataPBR{}) - ShaderBaseDataStart)
}

func NewShaderDataBase() ShaderDataBase {
	sdb := ShaderDataBase{}
	sdb.Setup()
	return sdb
}

func (s *ShaderDataBase) Setup() {
	s.SetModel(matrix.Mat4Identity())
}

func (s *ShaderDataBase) Destroy()           { s.destroyed = true }
func (s *ShaderDataBase) CancelDestroy()     { s.destroyed = false }
func (s *ShaderDataBase) IsDestroyed() bool  { return s.destroyed }
func (s *ShaderDataBase) IsActive() bool     { return !s.deactivated }
func (s *ShaderDataBase) Model() matrix.Mat4 { return s.model }

func (s *ShaderDataBase) Activate() {
	s.deactivated = false
	if s.shadow != nil {
		s.shadow.Activate()
	}
}

func (s *ShaderDataBase) Deactivate() {
	s.deactivated = true
	if s.shadow != nil {
		s.shadow.Deactivate()
	}
}

func (s *ShaderDataBase) setTransform(transform *matrix.Transform) {
	s.transform = transform
}

func (s *ShaderDataBase) setShadow(shadow DrawInstance) {
	s.shadow = shadow
	if s.deactivated {
		s.shadow.Deactivate()
	}
}

func (s *ShaderDataBase) SetModel(model matrix.Mat4) {
	s.InitModel = model
	if s.transform == nil {
		s.model = model
	}
}

func (s *ShaderDataBase) UpdateModel() {
	if s.transform != nil && s.transform.IsDirty() {
		s.model = matrix.Mat4Multiply(s.InitModel, s.transform.WorldMatrix())
	}
}

func (s *ShaderDataBase) DataPointer() unsafe.Pointer {
	return unsafe.Pointer(&s.model[0])
}

func (s *ShaderDataBase) UpdateNamedData(index, capacity int, name string) bool { return false }

func (s *ShaderDataBase) NamedDataPointer(name string) unsafe.Pointer { return nil }

func (s *ShaderDataBase) NamedDataInstanceSize(name string) int { return 0 }

type InstanceCopyData struct {
	bytes   []byte
	padding int
}

func InstanceCopyDataNew(padding int) InstanceCopyData {
	return InstanceCopyData{
		bytes:   make([]byte, 0),
		padding: padding,
	}
}

type DrawInstanceGroup struct {
	Mesh *Mesh
	InstanceDriverData
	MaterialInstance  *Material
	Instances         []DrawInstance
	rawData           InstanceCopyData
	namedInstanceData map[string]InstanceCopyData
	instanceSize      int
	visibleCount      int
	sort              int
	useBlending       bool
	destroyed         bool
}

func NewDrawInstanceGroup(mesh *Mesh, dataSize int) DrawInstanceGroup {
	return DrawInstanceGroup{
		Mesh:              mesh,
		Instances:         make([]DrawInstance, 0),
		rawData:           InstanceCopyDataNew(dataSize % 16),
		namedInstanceData: make(map[string]InstanceCopyData),
		instanceSize:      dataSize,
		destroyed:         false,
	}
}

func (d *DrawInstanceGroup) AlterPadding(blockSize int) {
	newPadding := blockSize - d.instanceSize%blockSize
	if d.rawData.padding != newPadding {
		d.rawData.padding = newPadding
		old := d.rawData.bytes
		d.rawData.bytes = make([]byte, d.TotalSize())
		copy(d.rawData.bytes, old)
	}
}

func (d *DrawInstanceGroup) IsEmpty() bool {
	return len(d.Instances) == 0
}

func (d *DrawInstanceGroup) IsReady() bool {
	// TODO:  Check if textures are ready?
	return d.Mesh.IsReady() && !d.IsEmpty()
}

func (d *DrawInstanceGroup) TotalSize() int {
	return len(d.Instances) * (d.instanceSize + d.rawData.padding)
}

func (d *DrawInstanceGroup) AddInstance(instance DrawInstance) {
	d.Instances = append(d.Instances, instance)
	d.rawData.bytes = append(d.rawData.bytes, make([]byte, d.instanceSize+d.rawData.padding)...)
	for i := range d.MaterialInstance.shaderInfo.LayoutGroups {
		g := &d.MaterialInstance.shaderInfo.LayoutGroups[i]
		for j := range g.Layouts {
			if g.Layouts[j].IsBuffer() {
				b := &g.Layouts[j]
				n := b.FullName()
				s := d.namedInstanceData[n]
				if len(s.bytes) < b.Capacity() {
					s.bytes = append(s.bytes, make([]byte, instance.NamedDataInstanceSize(n)+s.padding)...)
					d.namedInstanceData[n] = s
				}
			}
		}
	}
}

func (d *DrawInstanceGroup) texSize() (int32, int32) {
	// Low end devices have a max 2048 texture size
	pixelCount := int32(len(d.rawData.bytes)) / 4 / 4
	width := min(pixelCount, 2048)
	height := int32(1)
	for pixelCount > 2048 {
		height++
		pixelCount -= 2048
	}
	if height > 2048 {
		// TODO:  Handle this case with multiple textures
		panic("Too many instances")
	}
	return width, height
}

func (d *DrawInstanceGroup) AnyVisible() bool  { return d.visibleCount > 0 }
func (d *DrawInstanceGroup) VisibleCount() int { return d.visibleCount }

func (d *DrawInstanceGroup) VisibleSize() int {
	return d.visibleCount * (d.instanceSize + d.rawData.padding)
}

func (d *DrawInstanceGroup) updateNamedData(index int, instance DrawInstance, name string) {
	if !instance.UpdateNamedData(index, d.namedBuffers[name].capacity, name) {
		return
	}
	if ptr := instance.NamedDataPointer(name); ptr != nil {
		offset := uintptr(d.namedBuffers[name].stride * index)
		base := unsafe.Pointer(&d.namedInstanceData[name].bytes[0])
		to := unsafe.Pointer(uintptr(base) + offset)
		klib.Memcpy(to, ptr, uint64(len(d.namedInstanceData[name].bytes)))
	}
}

func (d *DrawInstanceGroup) UpdateData(renderer Renderer) {
	base := unsafe.Pointer(&d.rawData.bytes[0])
	offset := uintptr(0)
	count := len(d.Instances)
	d.visibleCount = 0
	instanceIndex := 0
	for i := 0; i < count; i++ {
		instance := d.Instances[i]
		instance.UpdateModel()
		if instance.IsDestroyed() {
			d.Instances[i] = d.Instances[count-1]
			i--
			count--
		} else if instance.IsActive() {
			if d.generatedSets {
				for k := range d.namedInstanceData {
					d.updateNamedData(instanceIndex, instance, k)
				}
			}
			to := unsafe.Pointer(uintptr(base) + offset)
			klib.Memcpy(to, instance.DataPointer(), uint64(d.instanceSize))
			offset += uintptr(d.instanceSize + d.rawData.padding)
			d.visibleCount++
			instanceIndex++
		}
	}
	if count < len(d.Instances) {
		newMemLen := count * (d.instanceSize + d.rawData.padding)
		d.Instances = d.Instances[:count]
		d.rawData.bytes = d.rawData.bytes[:newMemLen]
	}
	d.bindInstanceDriverData()
	if len(d.Instances) == 0 {
		renderer.DestroyGroup(d)
		d.destroyed = true
	}
}

func (d *DrawInstanceGroup) Clear(renderer Renderer) {
	if d.destroyed {
		return
	}
	for i := range d.Instances {
		d.Instances[i].Destroy()
	}
}

func (d *DrawInstanceGroup) Destroy(renderer Renderer) {
	if d.destroyed {
		return
	}
	d.Clear(renderer)
	d.Instances = d.Instances[:0]
	renderer.DestroyGroup(d)
}
