/******************************************************************************/
/* css_text_align.go                                                          */
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

package properties

import (
	"fmt"
	"kaiju/engine"
	"kaiju/engine/ui"
	"kaiju/engine/ui/markup/css/rules"
	"kaiju/engine/ui/markup/document"
	"kaiju/rendering"
)

// left|right|center|justify|initial|inherit
func (p TextAlign) Process(panel *ui.Panel, elm *document.Element, values []rules.PropertyValue, host *engine.Host) error {
	if len(values) != 1 {
		return fmt.Errorf("expected exactly 1 value but got %d", len(values))
	}
	labels := childLabels(elm)
	switch values[0].Str {
	case "left":
		for _, l := range labels {
			base := l.Base()
			base.Layout().AnchorTo(base.Layout().Anchor().ConvertToLeft())
			l.SetJustify(rendering.FontJustifyLeft)
		}
	case "right":
		for _, l := range labels {
			base := l.Base()
			base.Layout().AnchorTo(base.Layout().Anchor().ConvertToRight())
			l.SetJustify(rendering.FontJustifyRight)
		}
	case "center":
		for _, l := range labels {
			base := l.Base()
			base.Layout().AnchorTo(base.Layout().Anchor().ConvertToCenter())
			l.SetJustify(rendering.FontJustifyCenter)
		}
	case "justify":
		// TODO:  Support text justification
	case "initial":
	case "inherit":
	}
	return nil
}
