/******************************************************************************/
/* shader_pipeline.go                                                         */
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

package rendering

import (
	"log/slog"

	"kaijuengine.com/klib"
	vk "kaijuengine.com/rendering/vulkan"
	"kaijuengine.com/rendering/vulkan_const"
)

type ShaderPipelineData struct {
	Name                  string
	InputAssembly         ShaderPipelineInputAssembly
	Rasterization         ShaderPipelinePipelineRasterization
	Multisample           ShaderPipelinePipelineMultisample
	ColorBlend            ShaderPipelineColorBlend
	ColorBlendAttachments []ShaderPipelineColorBlendAttachments
	DepthStencil          ShaderPipelineDepthStencil
	Tessellation          ShaderPipelineTessellation
	GraphicsPipeline      ShaderPipelineGraphicsPipeline
	PushConstant          ShaderPipelinePushConstant
}

type ShaderPipelineInputAssembly struct {
	Topology         string `options:"StringVkPrimitiveTopology"`
	PrimitiveRestart bool
}

type ShaderPipelinePipelineRasterization struct {
	DepthClampEnable        bool
	RasterizerDiscardEnable bool
	PolygonMode             string `options:"StringVkPolygonMode"`
	CullMode                string `options:"StringVkCullModeFlagBits"`
	FrontFace               string `options:"StringVkFrontFace"`
	DepthBiasEnable         bool
	DepthBiasConstantFactor float32
	DepthBiasClamp          float32
	DepthBiasSlopeFactor    float32
	LineWidth               float32
}

type ShaderPipelinePipelineMultisample struct {
	RasterizationSamples  string `options:"StringVkSampleCountFlagBits"`
	SampleShadingEnable   bool
	MinSampleShading      float32
	AlphaToCoverageEnable bool
	AlphaToOneEnable      bool
}

type ShaderPipelineColorBlend struct {
	LogicOpEnable   bool
	LogicOp         string  `options:"StringVkLogicOp"`
	BlendConstants0 float32 `tip:"BlendConstants"`
	BlendConstants1 float32 `tip:"BlendConstants"`
	BlendConstants2 float32 `tip:"BlendConstants"`
	BlendConstants3 float32 `tip:"BlendConstants"`
}

type ShaderPipelineDepthStencil struct {
	DepthTestEnable       bool
	DepthWriteEnable      bool
	DepthCompareOp        string `options:"StringVkCompareOp"`
	DepthBoundsTestEnable bool
	StencilTestEnable     bool
	FrontFailOp           string `options:"StringVkStencilOp" tip:"FailOp"`
	FrontPassOp           string `options:"StringVkStencilOp" tip:"PassOp"`
	FrontDepthFailOp      string `options:"StringVkStencilOp" tip:"DepthFailOp"`
	FrontCompareOp        string `options:"StringVkCompareOp" tip:"CompareOp"`
	FrontCompareMask      uint32 `tip:"CompareMask"`
	FrontWriteMask        uint32 `tip:"WriteMask"`
	FrontReference        uint32 `tip:"Reference"`
	BackFailOp            string `options:"StringVkStencilOp" tip:"FailOp"`
	BackPassOp            string `options:"StringVkStencilOp" tip:"PassOp"`
	BackDepthFailOp       string `options:"StringVkStencilOp" tip:"DepthFailOp"`
	BackCompareOp         string `options:"StringVkCompareOp" tip:"CompareOp"`
	BackCompareMask       uint32 `tip:"CompareMask"`
	BackWriteMask         uint32 `tip:"WriteMask"`
	BackReference         uint32 `tip:"Reference"`
	MinDepthBounds        float32
	MaxDepthBounds        float32
}

type ShaderPipelineTessellation struct {
	PatchControlPoints string `options:"StringVkPatchControlPoints"`
}

type ShaderPipelineGraphicsPipeline struct {
	Subpass             uint32
	PipelineCreateFlags []string `options:"StringVkPipelineCreateFlagBits"`
}

type ShaderPipelinePushConstant struct {
	Size       uint32
	StageFlags []string `options:"StringVkAccessFlagBits"`
}

type ShaderPipelineColorBlendAttachments struct {
	BlendEnable         bool
	SrcColorBlendFactor string   `options:"StringVkBlendFactor"`
	DstColorBlendFactor string   `options:"StringVkBlendFactor"`
	ColorBlendOp        string   `options:"StringVkBlendOp"`
	SrcAlphaBlendFactor string   `options:"StringVkBlendFactor"`
	DstAlphaBlendFactor string   `options:"StringVkBlendFactor"`
	AlphaBlendOp        string   `options:"StringVkBlendOp"`
	ColorWriteMask      []string `options:"StringVkColorComponentFlagBits"`
}

type ShaderPipelineDataCompiled struct {
	Name                  string
	InputAssembly         ShaderPipelineInputAssemblyCompiled
	Rasterization         ShaderPipelinePipelineRasterizationCompiled
	Multisample           ShaderPipelinePipelineMultisampleCompiled
	ColorBlend            ShaderPipelineColorBlendCompiled
	ColorBlendAttachments []ShaderPipelineColorBlendAttachmentsCompiled
	DepthStencil          ShaderPipelineDepthStencilCompiled
	Tessellation          ShaderPipelineTessellationCompiled
	GraphicsPipeline      ShaderPipelineGraphicsPipelineCompiled
	PushConstant          ShaderPipelinePushConstantCompiled
}

type ShaderPipelineInputAssemblyCompiled struct {
	Topology         vulkan_const.PrimitiveTopology
	PrimitiveRestart bool
}

type ShaderPipelinePipelineRasterizationCompiled struct {
	DepthClampEnable        bool
	DiscardEnable           bool
	PolygonMode             vulkan_const.PolygonMode
	CullMode                vk.CullModeFlags
	FrontFace               vulkan_const.FrontFace
	DepthBiasEnable         bool
	DepthBiasConstantFactor float32
	DepthBiasClamp          float32
	DepthBiasSlopeFactor    float32
	LineWidth               float32
}

type ShaderPipelinePipelineMultisampleCompiled struct {
	RasterizationSamples  vulkan_const.SampleCountFlagBits
	SampleShadingEnable   bool
	MinSampleShading      float32
	AlphaToCoverageEnable bool
	AlphaToOneEnable      bool
}

type ShaderPipelineColorBlendCompiled struct {
	LogicOpEnable  bool
	LogicOp        vulkan_const.LogicOp
	BlendConstants [4]float32
}

type ShaderPipelineDepthStencilCompiled struct {
	DepthTestEnable       bool
	DepthWriteEnable      bool
	DepthCompareOp        vulkan_const.CompareOp
	DepthBoundsTestEnable bool
	StencilTestEnable     bool
	Front                 vk.StencilOpState
	Back                  vk.StencilOpState
	MinDepthBounds        float32
	MaxDepthBounds        float32
}

type ShaderPipelineTessellationCompiled struct {
	PatchControlPoints uint32
}

type ShaderPipelineGraphicsPipelineCompiled struct {
	Subpass             uint32
	PipelineCreateFlags vk.PipelineCreateFlags
}

type ShaderPipelinePushConstantCompiled struct {
	Size       uint32
	StageFlags vk.ShaderStageFlags
}

type ShaderPipelineColorBlendAttachmentsCompiled struct {
	BlendEnable         bool
	SrcColorBlendFactor vulkan_const.BlendFactor
	DstColorBlendFactor vulkan_const.BlendFactor
	ColorBlendOp        vulkan_const.BlendOp
	SrcAlphaBlendFactor vulkan_const.BlendFactor
	DstAlphaBlendFactor vulkan_const.BlendFactor
	AlphaBlendOp        vulkan_const.BlendOp
	ColorWriteMask      vk.ColorComponentFlags
}

func (d *ShaderPipelineData) Compile(device *GPUPhysicalDevice) ShaderPipelineDataCompiled {
	c := ShaderPipelineDataCompiled{
		Name: d.Name,
		InputAssembly: ShaderPipelineInputAssemblyCompiled{
			Topology:         d.InputAssembly.TopologyToVK(),
			PrimitiveRestart: d.InputAssembly.PrimitiveRestart,
		},
		Rasterization: ShaderPipelinePipelineRasterizationCompiled{
			DepthClampEnable:        d.Rasterization.DepthClampEnable,
			DiscardEnable:           d.Rasterization.RasterizerDiscardEnable,
			PolygonMode:             d.Rasterization.PolygonModeToVK(),
			CullMode:                vk.CullModeFlags(d.Rasterization.CullModeToVK()),
			FrontFace:               d.Rasterization.FrontFaceToVK(),
			DepthBiasEnable:         d.Rasterization.DepthBiasEnable,
			DepthBiasConstantFactor: d.Rasterization.DepthBiasConstantFactor,
			DepthBiasClamp:          d.Rasterization.DepthBiasClamp,
			DepthBiasSlopeFactor:    d.Rasterization.DepthBiasSlopeFactor,
			LineWidth:               d.Rasterization.LineWidth,
		},
		Multisample: ShaderPipelinePipelineMultisampleCompiled{
			RasterizationSamples:  vulkan_const.SampleCountFlagBits(d.Multisample.RasterizationSamplesToVK(device).toVulkan()),
			SampleShadingEnable:   d.Multisample.SampleShadingEnable,
			MinSampleShading:      d.Multisample.MinSampleShading,
			AlphaToCoverageEnable: d.Multisample.AlphaToCoverageEnable,
			AlphaToOneEnable:      d.Multisample.AlphaToOneEnable,
		},
		ColorBlend: ShaderPipelineColorBlendCompiled{
			LogicOpEnable: d.ColorBlend.LogicOpEnable,
			LogicOp:       d.ColorBlend.LogicOpToVK(),
			BlendConstants: [4]float32{
				d.ColorBlend.BlendConstants0,
				d.ColorBlend.BlendConstants1,
				d.ColorBlend.BlendConstants2,
				d.ColorBlend.BlendConstants3,
			},
		},
		ColorBlendAttachments: make([]ShaderPipelineColorBlendAttachmentsCompiled, len(d.ColorBlendAttachments)),
		DepthStencil: ShaderPipelineDepthStencilCompiled{
			DepthTestEnable:       d.DepthStencil.DepthTestEnable,
			DepthWriteEnable:      d.DepthStencil.DepthWriteEnable,
			DepthCompareOp:        compareOpToVK(d.DepthStencil.DepthCompareOp),
			DepthBoundsTestEnable: d.DepthStencil.DepthBoundsTestEnable,
			StencilTestEnable:     d.DepthStencil.StencilTestEnable,
			Front: vk.StencilOpState{
				FailOp:      stencilOpToVK(d.DepthStencil.FrontFailOp),
				PassOp:      stencilOpToVK(d.DepthStencil.FrontPassOp),
				DepthFailOp: stencilOpToVK(d.DepthStencil.FrontDepthFailOp),
				CompareOp:   compareOpToVK(d.DepthStencil.FrontCompareOp),
				CompareMask: d.DepthStencil.FrontCompareMask,
				WriteMask:   d.DepthStencil.FrontWriteMask,
				Reference:   d.DepthStencil.FrontReference,
			},
			Back: vk.StencilOpState{
				FailOp:      stencilOpToVK(d.DepthStencil.BackFailOp),
				PassOp:      stencilOpToVK(d.DepthStencil.BackPassOp),
				DepthFailOp: stencilOpToVK(d.DepthStencil.BackDepthFailOp),
				CompareOp:   compareOpToVK(d.DepthStencil.BackCompareOp),
				CompareMask: d.DepthStencil.BackCompareMask,
				WriteMask:   d.DepthStencil.BackWriteMask,
				Reference:   d.DepthStencil.BackReference,
			},
			MinDepthBounds: d.DepthStencil.MinDepthBounds,
			MaxDepthBounds: d.DepthStencil.MaxDepthBounds,
		},
		Tessellation: ShaderPipelineTessellationCompiled{
			PatchControlPoints: d.Tessellation.PatchControlPointsToVK(),
		},
		GraphicsPipeline: ShaderPipelineGraphicsPipelineCompiled{
			Subpass:             d.GraphicsPipeline.Subpass,
			PipelineCreateFlags: d.GraphicsPipeline.PipelineCreateFlagsToVK(),
		},
		PushConstant: ShaderPipelinePushConstantCompiled{
			Size:       d.PushConstant.Size,
			StageFlags: d.PushConstant.ShaderStageFlagsToVK(),
		},
	}
	for i := range d.ColorBlendAttachments {
		from := &d.ColorBlendAttachments[i]
		c.ColorBlendAttachments[i] = ShaderPipelineColorBlendAttachmentsCompiled{
			BlendEnable:         from.BlendEnable,
			SrcColorBlendFactor: from.SrcColorBlendFactorToVK(),
			DstColorBlendFactor: from.DstColorBlendFactorToVK(),
			ColorBlendOp:        from.ColorBlendOpToVK(),
			SrcAlphaBlendFactor: from.SrcAlphaBlendFactorToVK(),
			DstAlphaBlendFactor: from.DstAlphaBlendFactorToVK(),
			AlphaBlendOp:        from.AlphaBlendOpToVK(),
			ColorWriteMask:      vk.ColorComponentFlags(from.ColorWriteMaskToVK()),
		}
	}
	return c
}

func (a *ShaderPipelineColorBlendAttachments) ListSrcColorBlendFactor() []string {
	return klib.MapKeysSorted(StringVkBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) ListDstColorBlendFactor() []string {
	return klib.MapKeysSorted(StringVkBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) ListColorBlendOp() []string {
	return klib.MapKeysSorted(StringVkBlendOp)
}

func (a *ShaderPipelineColorBlendAttachments) ListSrcAlphaBlendFactor() []string {
	return klib.MapKeysSorted(StringVkBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) ListDstAlphaBlendFactor() []string {
	return klib.MapKeysSorted(StringVkBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) ListAlphaBlendOp() []string {
	return klib.MapKeysSorted(StringVkBlendOp)
}

func (a *ShaderPipelineColorBlendAttachments) SrcColorBlendFactorToVK() vulkan_const.BlendFactor {
	return blendFactorToVK(a.SrcColorBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) DstColorBlendFactorToVK() vulkan_const.BlendFactor {
	return blendFactorToVK(a.DstColorBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) ColorBlendOpToVK() vulkan_const.BlendOp {
	return blendOpToVK(a.ColorBlendOp)
}

func (a *ShaderPipelineColorBlendAttachments) SrcAlphaBlendFactorToVK() vulkan_const.BlendFactor {
	return blendFactorToVK(a.SrcAlphaBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) DstAlphaBlendFactorToVK() vulkan_const.BlendFactor {
	return blendFactorToVK(a.DstAlphaBlendFactor)
}

func (a *ShaderPipelineColorBlendAttachments) AlphaBlendOpToVK() vulkan_const.BlendOp {
	return blendOpToVK(a.AlphaBlendOp)
}

func (a *ShaderPipelineColorBlendAttachments) ColorWriteMaskToVK() vulkan_const.ColorComponentFlagBits {
	mask := vulkan_const.ColorComponentFlagBits(0)
	for i := range a.ColorWriteMask {
		mask |= StringVkColorComponentFlagBits[a.ColorWriteMask[i]]
	}
	return mask
}

func (s ShaderPipelineData) ListTopology() []string {
	return klib.MapKeysSorted(StringVkPrimitiveTopology)
}

func (s ShaderPipelineData) ListPolygonMode() []string {
	return klib.MapKeysSorted(StringVkPolygonMode)
}

func (s ShaderPipelineData) ListCullMode() []string {
	return klib.MapKeysSorted(StringVkCullModeFlagBits)
}

func (s ShaderPipelineData) ListFrontFace() []string {
	return klib.MapKeysSorted(StringVkFrontFace)
}

func (s ShaderPipelineData) ListRasterizationSamples() []string {
	return klib.MapKeysSorted(StringVkSampleCountFlagBits)
}

func (s ShaderPipelineData) ListBlendFactor() []string {
	return klib.MapKeysSorted(StringVkBlendFactor)
}

func (s ShaderPipelineData) ListBlendOp() []string {
	return klib.MapKeysSorted(StringVkBlendOp)
}

func (s ShaderPipelineData) ListLogicOp() []string {
	return klib.MapKeysSorted(StringVkLogicOp)
}

func (s ShaderPipelineData) ListDepthCompareOp() []string {
	return klib.MapKeysSorted(StringVkCompareOp)
}

func (s ShaderPipelineData) ListBackCompareOp() []string {
	return klib.MapKeysSorted(StringVkCompareOp)
}

func (s ShaderPipelineData) ListFrontFailOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListFrontPassOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListFrontDepthFailOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListFrontCompareOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListBackFailOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListBackPassOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListBackDepthFailOp() []string {
	return klib.MapKeysSorted(StringVkStencilOp)
}

func (s ShaderPipelineData) ListPatchControlPoints() []string {
	return klib.MapKeysSorted(StringVkPatchControlPoints)
}

func (s *ShaderPipelineData) PrimitiveRestartToVK() vk.Bool32 {
	return boolToVkBool(s.InputAssembly.PrimitiveRestart)
}

func (s *ShaderPipelineData) DepthClampEnableToVK() vk.Bool32 {
	return boolToVkBool(s.Rasterization.DepthClampEnable)
}

func (s *ShaderPipelineData) RasterizerDiscardEnableToVK() vk.Bool32 {
	return boolToVkBool(s.Rasterization.RasterizerDiscardEnable)
}

func (s *ShaderPipelineData) DepthBiasEnableToVK() vk.Bool32 {
	return boolToVkBool(s.Rasterization.DepthBiasEnable)
}

func (s *ShaderPipelineData) SampleShadingEnableToVK() vk.Bool32 {
	return boolToVkBool(s.Multisample.SampleShadingEnable)
}

func (s *ShaderPipelineData) AlphaToCoverageEnableToVK() vk.Bool32 {
	return boolToVkBool(s.Multisample.AlphaToCoverageEnable)
}

func (s *ShaderPipelineData) AlphaToOneEnableToVK() vk.Bool32 {
	return boolToVkBool(s.Multisample.AlphaToOneEnable)
}

func (s *ShaderPipelineData) LogicOpEnableToVK() vk.Bool32 {
	return boolToVkBool(s.ColorBlend.LogicOpEnable)
}

func (s *ShaderPipelineData) DepthTestEnableToVK() vk.Bool32 {
	return boolToVkBool(s.DepthStencil.DepthTestEnable)
}

func (s *ShaderPipelineData) DepthWriteEnableToVK() vk.Bool32 {
	return boolToVkBool(s.DepthStencil.DepthWriteEnable)
}

func (s *ShaderPipelineData) DepthBoundsTestEnableToVK() vk.Bool32 {
	return boolToVkBool(s.DepthStencil.DepthBoundsTestEnable)
}

func (s *ShaderPipelineData) StencilTestEnableToVK() vk.Bool32 {
	return boolToVkBool(s.DepthStencil.StencilTestEnable)
}

func (s *ShaderPipelineInputAssembly) TopologyToVK() vulkan_const.PrimitiveTopology {
	if res, ok := StringVkPrimitiveTopology[s.Topology]; ok {
		return res
	} else if s.Topology != "" {
		slog.Warn("invalid string for vkPrimitiveTopology", "value", s.Topology)
	}
	return vulkan_const.PrimitiveTopologyTriangleList
}

func (s *ShaderPipelinePipelineRasterization) PolygonModeToVK() vulkan_const.PolygonMode {
	if res, ok := StringVkPolygonMode[s.PolygonMode]; ok {
		return res
	} else if s.PolygonMode != "" {
		slog.Warn("invalid string for vkPolygonMode", "value", s.PolygonMode)
	}
	return vulkan_const.PolygonModeFill
}

func (s *ShaderPipelinePipelineRasterization) CullModeToVK() vulkan_const.CullModeFlagBits {
	if res, ok := StringVkCullModeFlagBits[s.CullMode]; ok {
		return res
	} else if s.CullMode != "" {
		slog.Warn("invalid string for vkCullModeFlagBits", "value", s.CullMode)
	}
	return vulkan_const.CullModeFrontBit
}

func (s *ShaderPipelinePipelineRasterization) FrontFaceToVK() vulkan_const.FrontFace {
	if res, ok := StringVkFrontFace[s.FrontFace]; ok {
		return res
	} else if s.FrontFace != "" {
		slog.Warn("invalid string for vkFrontFace", "value", s.FrontFace)
	}
	return vulkan_const.FrontFaceClockwise
}

func (s *ShaderPipelinePipelineMultisample) RasterizationSamplesToVK(device *GPUPhysicalDevice) GPUSampleCountFlags {
	return sampleCountToGpu(s.RasterizationSamples, device)
}

func (s *ShaderPipelineColorBlend) LogicOpToVK() vulkan_const.LogicOp {
	if res, ok := StringVkLogicOp[s.LogicOp]; ok {
		return res
	} else if s.LogicOp != "" {
		slog.Warn("invalid string for vkLogicOp", "value", s.LogicOp)
	}
	return vulkan_const.LogicOpCopy
}

func (s *ShaderPipelineData) BlendConstants() [4]float32 {
	return [4]float32{
		s.ColorBlend.BlendConstants0,
		s.ColorBlend.BlendConstants1,
		s.ColorBlend.BlendConstants2,
		s.ColorBlend.BlendConstants3,
	}
}

func (s *ShaderPipelineTessellation) PatchControlPointsToVK() uint32 {
	if res, ok := StringVkPatchControlPoints[s.PatchControlPoints]; ok {
		return res
	} else if s.PatchControlPoints != "" {
		slog.Warn("invalid string for PatchControlPoints", "value", s.PatchControlPoints)
	}
	return 3
}

// TODO:  This and the BackStencilOpStateToVK are duplicates because of a bad
// structure setup, please fix later
func (s *ShaderPipelineData) FrontStencilOpStateToVK() vk.StencilOpState {
	return vk.StencilOpState{
		FailOp:      stencilOpToVK(s.DepthStencil.FrontFailOp),
		PassOp:      stencilOpToVK(s.DepthStencil.FrontPassOp),
		DepthFailOp: stencilOpToVK(s.DepthStencil.FrontDepthFailOp),
		CompareOp:   compareOpToVK(s.DepthStencil.FrontCompareOp),
		CompareMask: s.DepthStencil.FrontCompareMask,
		WriteMask:   s.DepthStencil.FrontWriteMask,
		Reference:   s.DepthStencil.FrontReference,
	}
}

func (s *ShaderPipelineData) BackStencilOpStateToVK() vk.StencilOpState {
	return vk.StencilOpState{
		FailOp:      stencilOpToVK(s.DepthStencil.BackFailOp),
		PassOp:      stencilOpToVK(s.DepthStencil.BackPassOp),
		DepthFailOp: stencilOpToVK(s.DepthStencil.BackDepthFailOp),
		CompareOp:   compareOpToVK(s.DepthStencil.BackCompareOp),
		CompareMask: s.DepthStencil.BackCompareMask,
		WriteMask:   s.DepthStencil.BackWriteMask,
		Reference:   s.DepthStencil.BackReference,
	}
}

func (s *ShaderPipelineGraphicsPipeline) PipelineCreateFlagsToVK() vk.PipelineCreateFlags {
	mask := vulkan_const.PipelineCreateFlagBits(0)
	for i := range s.PipelineCreateFlags {
		mask |= StringVkPipelineCreateFlagBits[s.PipelineCreateFlags[i]]
	}
	return vk.PipelineCreateFlags(mask)
}

func (s *ShaderPipelinePushConstant) ShaderStageFlagsToVK() vk.ShaderStageFlags {
	mask := vulkan_const.ShaderStageFlagBits(0)
	for i := range s.StageFlags {
		mask |= StringVkShaderStageFlagBits[s.StageFlags[i]]
	}
	return vk.ShaderStageFlags(mask)
}
