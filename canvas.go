package hashcatlauncher

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/canvas"
)

func banner() fyne.CanvasObject {
	t := canvas.NewText("hashcat.launcher v"+Version, &color.RGBA{0xff, 0xff, 0xff, 0xff})
	t.TextSize = fyne.CurrentApp().Settings().Theme().Size("text") * 2
	r := canvas.NewRectangle(&color.RGBA{0, 0, 0, 0xff})
	return container.NewMax(
		r,
		container.NewCenter(
			t,
		),
	)
}
