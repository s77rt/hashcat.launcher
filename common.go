package hashcatlauncher

import (
	"fmt"
	"strings"
	"regexp"
	"image/color"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
var re_ansi = regexp.MustCompile(ansi)
var re_status = regexp.MustCompile(`^(Status|Hash\.Name|Hash\.Type|Hash\.Target|Guess\.Queue|Guess\.Base|Guess\.Mod|Guess\.Mask|Guess\.Charset|Guess\.Queue\.Base|Guess\.Queue\.Mod|Progress|Recovered|Time\.Started|Time\.Estimated|Speed\.#1|Speed.#\*)\.+:\s(.+)$`)
var re_progress = regexp.MustCompile(`(\d+\/\d+)\s\(([\d\.]+)\%\)`)
var re_recovered = regexp.MustCompile(`(\d+\/\d+)\s\(([\d\.]+)\%\)\sDigests`)
var re_speed = regexp.MustCompile(`[\d\.]+\s\wH\/s`)
var re_mode = regexp.MustCompile(`^MODE:\s(\d+)$`)
var re_type = regexp.MustCompile(`^TYPE:\s(.+)$`)

func spacer(w int, h int) fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{0, 0, 0, 0})
	rect.SetMinSize(fyne.Size{w, h})
	return rect
}

func reserved(w int, h int) fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{12, 5, 77, 0xff})
	rect.SetMinSize(fyne.Size{w, h})
	return rect
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
func intSliceToString(a []int, delim string) string {
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
    //return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
    //return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

