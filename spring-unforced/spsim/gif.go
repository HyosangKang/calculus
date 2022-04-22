package spsim

type Ginfo struct {
	Height   int    // height (pixel) of image
	Width    int    // width (pixel) of image
	SpWidth  int    // spring canvas width (pixel) -- on the left
	GphWidth int    // graph canvas width (pixel) -- on the right
	Pad      int    // padding size (pixel)
	Ceil     int    // ceiling height (pixel)
	SpL      int    // spring minimum length (pixel)
	CR       int    // circle radius (pixel)
	Coil     [3]int // coil info: coil number (int), coil radius (pixel), coil margin (pixel)
	Delay    int
}

func initGI() Ginfo {
	g := Ginfo{
		SpWidth:  200,
		GphWidth: 600,
		Ceil:     50,
		SpL:      100,
		CR:       10,
		Pad:      10,
		Coil:     [3]int{10, 10, 5},
		Delay:    1,
	}
	g.Height = g.Ceil + g.SpL + g.GphWidth
	g.Width = g.SpWidth + g.GphWidth
	return g
}
