/******************************************************************************/
/* css_background_color.go                                                    */
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
	"kaiju/engine/assets"
	"kaiju/engine/ui"
	"kaiju/engine/ui/markup/css/functions"
	"kaiju/engine/ui/markup/css/helpers"
	"kaiju/engine/ui/markup/css/rules"
	"kaiju/engine/ui/markup/document"
	"kaiju/matrix"
)

func setChildTextBackgroundColor(elm *document.Element, color matrix.Color) {
	for _, c := range elm.Children {
		if c.IsText() {
			c.UI.ToLabel().SetBGColor(color)
		} else if c.UI.ToPanel().Background() == nil { // Only continue if transparent
			setChildTextBackgroundColor(c, color)
		}
	}
}

func (p BackgroundColor) Process(panel *ui.Panel, elm *document.Element, values []rules.PropertyValue, host *engine.Host) error {
	if len(values) != 1 {
		return fmt.Errorf("Expected exactly 1 value but got %d", len(values))
	}
	// Images used for background are not colored
	bg := elm.UI.ToPanel().Background()
	applyPanelColor := bg == nil || bg.Key == assets.TextureSquare
	var err error
	var color matrix.Color
	hex := values[0].Str
	if hex == "inherit" {
		if applyPanelColor {
			pBase := panel.Base()
			pBase.AddEvent(ui.EventTypeRender, func() {
				if pBase.Entity().Parent != nil {
					p := ui.FirstPanelOnEntity(pBase.Entity().Parent)
					panel.SetColor(p.Base().ShaderData().FgColor)
				}
			})
		}
		return nil
	} else {
		switch values[0].Str {
		case "rgb":
			hex, _ = functions.Rgb{}.Process(panel, elm, values[0])
		case "rgba":
			hex, _ = functions.Rgba{}.Process(panel, elm, values[0])
		}
		if newHex, ok := helpers.ColorMap[hex]; ok {
			hex = newHex
		}
		if color, err = matrix.ColorFromHexString(hex); err == nil {
			if applyPanelColor || panel.Base().Type() == ui.ElementTypeImage {
				panel.SetColor(color)
			}
			if !panel.HasEnforcedColor() {
				setChildTextBackgroundColor(elm, color)
			}
			return nil
		} else {
			return err
		}
	}
}
