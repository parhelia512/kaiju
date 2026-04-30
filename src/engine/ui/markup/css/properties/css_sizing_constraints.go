/******************************************************************************/
/* css_sizing_constraints.go                                                  */
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

package properties

import (
	"strconv"
	"strings"

	"kaijuengine.com/engine/ui"
)

const cssSizingConstraintsKey = "markup.css.sizing.constraints"

type cssSizingConstraints struct {
	MinWidth      float32
	MaxWidth      float32
	MinHeight     float32
	MaxHeight     float32
	AspectRatio   float32
	HasMinWidth   bool
	HasMaxWidth   bool
	HasMinHeight  bool
	HasMaxHeight  bool
	HasAspect     bool
	UseBorderBox  bool
	HasBoxSizing  bool
}

func currentSizingConstraints(panel *ui.Panel) cssSizingConstraints {
	data := panel.Base().Entity().NamedData(cssSizingConstraintsKey)
	if len(data) == 0 {
		return cssSizingConstraints{}
	}
	if c, ok := data[len(data)-1].(cssSizingConstraints); ok {
		return c
	}
	return cssSizingConstraints{}
}

func storeSizingConstraints(panel *ui.Panel, c cssSizingConstraints) {
	entity := panel.Base().Entity()
	entity.RemoveNamedDataByName(cssSizingConstraintsKey)
	entity.AddNamedData(cssSizingConstraintsKey, c)
}

func setMinWidth(panel *ui.Panel, v float32, enabled bool) {
	c := currentSizingConstraints(panel)
	c.MinWidth = v
	c.HasMinWidth = enabled
	storeSizingConstraints(panel, c)
}

func setMaxWidth(panel *ui.Panel, v float32, enabled bool) {
	c := currentSizingConstraints(panel)
	c.MaxWidth = v
	c.HasMaxWidth = enabled
	storeSizingConstraints(panel, c)
}

func setMinHeight(panel *ui.Panel, v float32, enabled bool) {
	c := currentSizingConstraints(panel)
	c.MinHeight = v
	c.HasMinHeight = enabled
	storeSizingConstraints(panel, c)
}

func setMaxHeight(panel *ui.Panel, v float32, enabled bool) {
	c := currentSizingConstraints(panel)
	c.MaxHeight = v
	c.HasMaxHeight = enabled
	storeSizingConstraints(panel, c)
}

func setAspectRatio(panel *ui.Panel, ratio float32, enabled bool) {
	c := currentSizingConstraints(panel)
	c.AspectRatio = ratio
	c.HasAspect = enabled
	storeSizingConstraints(panel, c)
}

func setBoxSizing(panel *ui.Panel, useBorderBox bool) {
	c := currentSizingConstraints(panel)
	c.UseBorderBox = useBorderBox
	c.HasBoxSizing = true
	storeSizingConstraints(panel, c)
}

func applyWidthConstraints(panel *ui.Panel, width float32) float32 {
	c := currentSizingConstraints(panel)
	if c.HasMinWidth && width < c.MinWidth {
		width = c.MinWidth
	}
	if c.HasMaxWidth && width > c.MaxWidth {
		width = c.MaxWidth
	}
	return width
}

func applyHeightConstraints(panel *ui.Panel, height float32) float32 {
	c := currentSizingConstraints(panel)
	if c.HasMinHeight && height < c.MinHeight {
		height = c.MinHeight
	}
	if c.HasMaxHeight && height > c.MaxHeight {
		height = c.MaxHeight
	}
	return height
}

func parseRatio(values []string) (float32, bool) {
	if len(values) == 1 {
		r := strings.TrimSpace(values[0])
		if r == "auto" || r == "initial" {
			return 0, false
		}
		if strings.Contains(r, "/") {
			parts := strings.Split(r, "/")
			if len(parts) == 2 {
				left := strings.TrimSpace(parts[0])
				right := strings.TrimSpace(parts[1])
				if left != "" && right != "" {
					lv := parseSimpleFloat(left)
					rv := parseSimpleFloat(right)
					if rv > 0 {
						return lv / rv, true
					}
				}
			}
		}
	}
	if len(values) == 3 && values[1] == "/" {
		lv := parseSimpleFloat(values[0])
		rv := parseSimpleFloat(values[2])
		if rv > 0 {
			return lv / rv, true
		}
	}
	return 0, false
}

func parseSimpleFloat(v string) float32 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	out, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0
	}
	return float32(out)
}
