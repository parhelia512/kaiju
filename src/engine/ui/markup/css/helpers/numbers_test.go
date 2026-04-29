package helpers

import (
	"testing"

	"kaijuengine.com/rendering"
)

type testWindow struct {
	dpmm   float64
	width  int
	height int
}

func (w testWindow) DotsPerMillimeter() float64 { return w.dpmm }
func (w testWindow) Width() int                 { return w.width }
func (w testWindow) Height() int                { return w.height }

func TestNumFromLengthWithFont_Units(t *testing.T) {
	w := testWindow{dpmm: 2, width: 1000, height: 500}
	fontSize := float32(20)

	tests := []struct {
		name string
		in   string
		want float32
	}{
		{name: "percent", in: "75%", want: 0.75},
		{name: "px", in: "100px", want: 100},
		{name: "em", in: "3em", want: 60},
		{name: "ex", in: "6ex", want: 120},
		{name: "cm", in: "4cm", want: 80},
		{name: "mm", in: "40mm", want: 80},
		{name: "in", in: "1.5in", want: 76.2},
		{name: "pt", in: "72pt", want: 50.8},
		{name: "pc", in: "6pc", want: 50.8},
		{name: "rem", in: "2rem", want: 2 * rendering.DefaultFontEMSize},
		{name: "vw", in: "30vw", want: 300},
		{name: "vh", in: "10vh", want: 50},
		{name: "vmin", in: "20vmin", want: 100},
		{name: "vmax", in: "20vmax", want: 200},
		{name: "ch", in: "20ch", want: 200}, // 20 * (0.5 * 20)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NumFromLengthWithFont(tt.in, w, fontSize)
			if got != tt.want {
				t.Fatalf("NumFromLengthWithFont(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNumFromLength_DefaultFontContext(t *testing.T) {
	w := testWindow{dpmm: 2, width: 1000, height: 500}

	got := NumFromLength("2em", w)
	want := float32(2) * rendering.DefaultFontEMSize
	if got != want {
		t.Fatalf("NumFromLength(%q) = %v, want %v", "2em", got, want)
	}
}
