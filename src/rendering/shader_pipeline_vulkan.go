package rendering

import (
	"log/slog"
	"unsafe"

	"kaijuengine.com/platform/profiler/tracing"
	vk "kaijuengine.com/rendering/vulkan"
	"kaijuengine.com/rendering/vulkan_const"
)

func (s *ShaderPipelineDataCompiled) ConstructPipeline(device *GPUDevice, shader *Shader, renderPass *RenderPass, stages []vk.PipelineShaderStageCreateInfo) bool {
	defer tracing.NewRegion("ShaderPipelineDataCompiled.ConstructPipeline").End()
	pSetLayout := vk.DescriptorSetLayout(shader.RenderId.descriptorSetLayout.handle)
	pipelineLayoutInfo := vk.PipelineLayoutCreateInfo{
		SType:          vulkan_const.StructureTypePipelineLayoutCreateInfo,
		Flags:          0, // PipelineLayoutCreateFlags
		SetLayoutCount: 1,
		PSetLayouts:    &pSetLayout,
	}
	if s.PushConstant.Size > 0 {
		pushRanges := [1]vk.PushConstantRange{{
			StageFlags: s.PushConstant.StageFlags,
			Offset:     0,
			Size:       s.PushConstant.Size,
		}}
		pipelineLayoutInfo.PushConstantRangeCount = 1
		pipelineLayoutInfo.PPushConstantRanges = &pushRanges[0]
	}
	var pLayout vk.PipelineLayout
	if vk.CreatePipelineLayout(vk.Device(device.LogicalDevice.handle), &pipelineLayoutInfo, nil, &pLayout) != vulkan_const.Success {
		slog.Error("Failed to create pipeline layout")
		return false
	} else {
		device.LogicalDevice.dbg.track(unsafe.Pointer(pLayout))
	}
	shader.RenderId.pipelineLayout.handle = unsafe.Pointer(pLayout)
	bDesc := vertexGetBindingDescription(shader)
	bDescCount := uint32(len(bDesc))
	for i := uint32(1); i < bDescCount; i++ {
		bDesc[i].Stride = uint32(device.PhysicalDevice.PadBufferSize(uintptr(bDesc[i].Stride)))
	}
	aDesc := vertexGetAttributeDescription(shader)
	vertexInputInfo := vk.PipelineVertexInputStateCreateInfo{
		SType:                           vulkan_const.StructureTypePipelineVertexInputStateCreateInfo,
		VertexBindingDescriptionCount:   bDescCount,
		VertexAttributeDescriptionCount: uint32(len(aDesc)),
		PVertexBindingDescriptions:      &bDesc[0],
		PVertexAttributeDescriptions:    &aDesc[0],
	}
	inputAssembly := vk.PipelineInputAssemblyStateCreateInfo{
		SType:                  vulkan_const.StructureTypePipelineInputAssemblyStateCreateInfo,
		Flags:                  0, // PipelineInputAssemblyStateCreateFlags
		Topology:               s.InputAssembly.Topology,
		PrimitiveRestartEnable: boolToVkBool(s.InputAssembly.PrimitiveRestart),
	}
	sce := device.LogicalDevice.SwapChain.Extent
	viewport := vk.Viewport{
		X:        0.0,
		Y:        0.0,
		Width:    float32(sce.Width()),
		Height:   float32(sce.Height()),
		MinDepth: 0.0,
		MaxDepth: 1.0,
	}
	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: vk.Extent2D{
			Width:  uint32(sce.Width()),
			Height: uint32(sce.Height()),
		},
	}
	dynamicStates := []vulkan_const.DynamicState{
		vulkan_const.DynamicStateViewport,
		vulkan_const.DynamicStateScissor,
	}
	dynamicState := vk.PipelineDynamicStateCreateInfo{
		SType:             vulkan_const.StructureTypePipelineDynamicStateCreateInfo,
		DynamicStateCount: uint32(len(dynamicStates)),
		PDynamicStates:    &dynamicStates[0],
	}
	viewportState := vk.PipelineViewportStateCreateInfo{
		SType:         vulkan_const.StructureTypePipelineViewportStateCreateInfo,
		ViewportCount: 1,
		PViewports:    &viewport,
		ScissorCount:  1,
		PScissors:     &scissor,
	}
	rasterizer := vk.PipelineRasterizationStateCreateInfo{
		SType:                   vulkan_const.StructureTypePipelineRasterizationStateCreateInfo,
		Flags:                   0, // PipelineRasterizationStateCreateFlags
		DepthClampEnable:        boolToVkBool(s.Rasterization.DepthClampEnable),
		RasterizerDiscardEnable: boolToVkBool(s.Rasterization.DiscardEnable),
		PolygonMode:             s.Rasterization.PolygonMode,
		LineWidth:               s.Rasterization.LineWidth,
		CullMode:                s.Rasterization.CullMode,
		FrontFace:               s.Rasterization.FrontFace,
		DepthBiasEnable:         boolToVkBool(s.Rasterization.DepthBiasEnable),
		DepthBiasConstantFactor: s.Rasterization.DepthBiasConstantFactor,
		DepthBiasClamp:          s.Rasterization.DepthBiasClamp,
		DepthBiasSlopeFactor:    s.Rasterization.DepthBiasSlopeFactor,
	}
	multisampling := vk.PipelineMultisampleStateCreateInfo{
		SType:                 vulkan_const.StructureTypePipelineMultisampleStateCreateInfo,
		Flags:                 0, // PipelineMultisampleStateCreateFlags
		SampleShadingEnable:   boolToVkBool(s.Multisample.SampleShadingEnable),
		RasterizationSamples:  s.Multisample.RasterizationSamples,
		MinSampleShading:      s.Multisample.MinSampleShading,
		PSampleMask:           nil,
		AlphaToCoverageEnable: boolToVkBool(s.Multisample.AlphaToCoverageEnable),
		AlphaToOneEnable:      boolToVkBool(s.Multisample.AlphaToOneEnable),
	}
	colorBlendAttachment := make([]vk.PipelineColorBlendAttachmentState, len(s.ColorBlendAttachments))
	for i := range s.ColorBlendAttachments {
		colorBlendAttachment[i].BlendEnable = boolToVkBool(s.ColorBlendAttachments[i].BlendEnable)
		colorBlendAttachment[i].SrcColorBlendFactor = s.ColorBlendAttachments[i].SrcColorBlendFactor
		colorBlendAttachment[i].DstColorBlendFactor = s.ColorBlendAttachments[i].DstColorBlendFactor
		colorBlendAttachment[i].ColorBlendOp = s.ColorBlendAttachments[i].ColorBlendOp
		colorBlendAttachment[i].SrcAlphaBlendFactor = s.ColorBlendAttachments[i].SrcAlphaBlendFactor
		colorBlendAttachment[i].DstAlphaBlendFactor = s.ColorBlendAttachments[i].DstAlphaBlendFactor
		colorBlendAttachment[i].AlphaBlendOp = s.ColorBlendAttachments[i].AlphaBlendOp
		writeMask := s.ColorBlendAttachments[i].ColorWriteMask
		colorBlendAttachment[i].ColorWriteMask = vk.ColorComponentFlags(writeMask)
	}
	colorBlendAttachmentCount := len(colorBlendAttachment)
	colorBlending := vk.PipelineColorBlendStateCreateInfo{
		SType:           vulkan_const.StructureTypePipelineColorBlendStateCreateInfo,
		Flags:           0, // PipelineColorBlendStateCreateFlags
		LogicOpEnable:   boolToVkBool(s.ColorBlend.LogicOpEnable),
		LogicOp:         s.ColorBlend.LogicOp,
		AttachmentCount: uint32(colorBlendAttachmentCount),
		BlendConstants:  s.ColorBlend.BlendConstants,
	}
	if colorBlendAttachmentCount > 0 {
		colorBlending.PAttachments = &colorBlendAttachment[0]
	}
	pipelineInfo := vk.GraphicsPipelineCreateInfo{
		SType:               vulkan_const.StructureTypeGraphicsPipelineCreateInfo,
		Flags:               s.GraphicsPipeline.PipelineCreateFlags,
		StageCount:          uint32(len(stages)),
		PStages:             &stages[0],
		PVertexInputState:   &vertexInputInfo,
		PInputAssemblyState: &inputAssembly,
		PViewportState:      &viewportState,
		PRasterizationState: &rasterizer,
		PMultisampleState:   &multisampling,
		PColorBlendState:    &colorBlending,
		PDynamicState:       &dynamicState,
		Layout:              vk.PipelineLayout(shader.RenderId.pipelineLayout.handle),
		RenderPass:          renderPass.Handle,
		BasePipelineHandle:  vk.Pipeline(vk.NullHandle),
		Subpass:             s.GraphicsPipeline.Subpass,
	}
	hasDepth := false
	for i := 0; i < len(renderPass.construction.SubpassDescriptions) && !hasDepth; i++ {
		hasDepth = len(renderPass.construction.SubpassDescriptions[i].DepthStencilAttachment) > 0
	}
	var depthStencil vk.PipelineDepthStencilStateCreateInfo
	if hasDepth {
		depthStencil = vk.PipelineDepthStencilStateCreateInfo{
			SType:                 vulkan_const.StructureTypePipelineDepthStencilStateCreateInfo,
			Flags:                 0, // PipelineDepthStencilStateCreateFlags
			DepthTestEnable:       boolToVkBool(s.DepthStencil.DepthTestEnable),
			DepthCompareOp:        s.DepthStencil.DepthCompareOp,
			DepthBoundsTestEnable: boolToVkBool(s.DepthStencil.DepthBoundsTestEnable),
			StencilTestEnable:     boolToVkBool(s.DepthStencil.StencilTestEnable),
			MinDepthBounds:        s.DepthStencil.MinDepthBounds,
			MaxDepthBounds:        s.DepthStencil.MaxDepthBounds,
			DepthWriteEnable:      boolToVkBool(s.DepthStencil.DepthWriteEnable),
			Front:                 s.DepthStencil.Front,
			Back:                  s.DepthStencil.Back,
		}
		pipelineInfo.PDepthStencilState = &depthStencil
	}
	tess := vk.PipelineTessellationStateCreateInfo{}
	if len(shader.data.TessellationControl) > 0 ||
		len(shader.data.TessellationEvaluation) > 0 {
		tess.SType = vulkan_const.StructureTypePipelineTessellationStateCreateInfo
		tess.Flags = 0 // PipelineTessellationStateCreateFlags
		tess.PatchControlPoints = s.Tessellation.PatchControlPoints
		pipelineInfo.PTessellationState = &tess
	}
	success := true
	pipelines := [1]vk.Pipeline{}
	if vk.CreateGraphicsPipelines(vk.Device(device.LogicalDevice.handle), vk.PipelineCache(vk.NullHandle), 1, &pipelineInfo, nil, &pipelines[0]) != vulkan_const.Success {
		success = false
		slog.Error("Failed to create graphics pipeline")
	} else {
		device.LogicalDevice.dbg.track(unsafe.Pointer(pipelines[0]))
	}
	shader.RenderId.graphicsPipeline.handle = unsafe.Pointer(pipelines[0])
	return success
}
