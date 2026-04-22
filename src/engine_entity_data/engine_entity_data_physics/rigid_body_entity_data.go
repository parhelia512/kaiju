/******************************************************************************/
/* rigid_body_data_binding.go                                                 */
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

package engine_entity_data_physics

import (
	"log/slog"

	"kaijuengine.com/engine"
	"kaijuengine.com/engine/encoding/pod"
	"kaijuengine.com/engine/physics"
	"kaijuengine.com/engine_entity_data/content_id"
	"kaijuengine.com/matrix"
	"kaijuengine.com/rendering/loaders/kaiju_mesh"
)

var bindingKey = ""

type Shape int

const (
	ShapeBox Shape = iota
	ShapeSphere
	ShapeCapsule
	ShapeCylinder
	ShapeCone
	ShapeMesh
)

func init() {
	engine.RegisterEntityData(RigidBodyEntityData{})
}

func BindingKey() string {
	if bindingKey == "" {
		bindingKey = pod.QualifiedNameForLayout(RigidBodyEntityData{})
	}
	return bindingKey
}

type RigidBodyEntityData struct {
	AssetKey content_id.Mesh
	Extent   matrix.Vec3 `default:"1,1,1"`
	Mass     float32     `default:"1"`
	Radius   float32     `default:"1"`
	Height   float32     `default:"1"`
	Shape    Shape
	IsStatic bool
}

func (r RigidBodyEntityData) Init(e *engine.Entity, host *engine.Host) {
	host.StartPhysics()
	t := &e.Transform
	scale := t.Scale()
	var shape *physics.CollisionShape
	switch r.Shape {
	case ShapeBox:
		size := r.Extent.Multiply(scale)
		shape = &physics.NewBoxShape(size).CollisionShape
	case ShapeSphere:
		rad := r.Radius * float32(scale.LongestAxisValue())
		shape = &physics.NewSphereShape(rad).CollisionShape
	case ShapeCapsule:
		rad := r.Radius * float32(scale.LongestAxisValue())
		height := r.Height * scale.Y()
		shape = &physics.NewCapsuleShape(rad, height).CollisionShape
	case ShapeCylinder:
		size := r.Extent.Multiply(scale)
		shape = &physics.NewCylinderShape(size).CollisionShape
	case ShapeCone:
		rad := r.Radius * float32(scale.LongestAxisValue())
		height := r.Height * scale.Y()
		shape = &physics.NewConeShape(rad, height).CollisionShape
	case ShapeMesh:
		data, err := host.AssetDatabase().Read(string(r.AssetKey))
		onErr := func() {
			slog.Error("Failed to read the asset for the physics mesh shape, falling back to a box", "error", err)
			size := r.Extent.Multiply(e.Transform.Scale())
			shape = &physics.NewBoxShape(size).CollisionShape
		}
		if err == nil {
			scale := e.Transform.WorldScale()
			km, err := kaiju_mesh.Deserialize(data)
			if err == nil {
				verts := make([]float32, len(km.Verts)*3)
				idx := 0
				for i := range km.Verts {
					// TODO:  This should probably use the transformation matrix
					// to also account for rotation, translation likely isn't needed
					pos := km.Verts[i].Position.Multiply(scale)
					verts[idx] = pos[matrix.Vx]
					verts[idx+1] = pos[matrix.Vy]
					verts[idx+2] = pos[matrix.Vz]
					idx += 3
				}
				triangleIVA := physics.NewTriangleIndexVertexArray(km.Indexes, verts)
				shape = &physics.NewBvhTriangleMeshShape(triangleIVA, false).CollisionShape
			} else {
				onErr()
			}
		} else {
			onErr()
		}
	}
	if r.IsStatic {
		r.Mass = 0
	}
	inertia := shape.CalculateLocalInertia(r.Mass)
	motion := physics.NewDefaultMotionState(matrix.QuaternionFromEuler(t.Rotation()), t.Position())
	body := physics.NewRigidBody(r.Mass, motion, shape, inertia)
	host.Physics().AddEntity(e, body)
}
