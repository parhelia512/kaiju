/******************************************************************************/
/* editor_config.go                                                           */
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

package editor_config

type FileExtension = string
type AssetType = string

const (
	FileExtensionH              FileExtension = ".h"
	FileExtensionC              FileExtension = ".c"
	FileExtensionGo             FileExtension = ".go"
	FileExtensionMap            FileExtension = ".map"
	FileExtensionObj            FileExtension = ".obj"
	FileExtensionGlb            FileExtension = ".glb"
	FileExtensionGltf           FileExtension = ".gltf"
	FileExtensionPng            FileExtension = ".png"
	FileExtensionMesh           FileExtension = ".msh"
	FileExtensionStage          FileExtension = ".stg"
	FileExtensionHTML           FileExtension = ".html"
	FileExtensionShader         FileExtension = ".shader"
	FileExtensionRenderPass     FileExtension = ".renderpass"
	FileExtensionShaderPipeline FileExtension = ".shaderpipeline"
	FileExtensionMaterial       FileExtension = ".material"
	FileExtensionAssetDbInfo    FileExtension = ".adi"
)

const (
	AssetTypeH              AssetType = "h"
	AssetTypeC              AssetType = "c"
	AssetTypeGo             AssetType = "go"
	AssetTypeMap            AssetType = "map"
	AssetTypeObj            AssetType = "obj"
	AssetTypeGlb            AssetType = "glb"
	AssetTypeGltf           AssetType = "gltf"
	AssetTypeImage          AssetType = "image"
	AssetTypeMesh           AssetType = "mesh"
	AssetTypeStage          AssetType = "stg"
	AssetTypeHTML           AssetType = "html"
	AssetTypeShader         AssetType = "shader"
	AssetTypeRenderPass     AssetType = "renderpass"
	AssetTypeShaderPipeline AssetType = "shaderpipeline"
	AssetTypeMaterial       AssetType = "material"
)
