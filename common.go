package hashcatlauncher

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
	"regexp"
	"image/color"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
var re_ansi = regexp.MustCompile(ansi)
var re_info = regexp.MustCompile(`INFO:\s.+potfile(?:\!|\.)`)
var re_status = regexp.MustCompile(`^(Status|Hash\.Name|Hash\.Type|Hash\.Target|Guess\.Queue|Guess\.Base|Guess\.Mod|Guess\.Mask|Guess\.Charset|Guess\.Queue\.Base|Guess\.Queue\.Mod|Progress|Recovered|Time\.Started|Time\.Estimated|Speed\.#1|Speed.#\*|Hardware\.Mon\.#\d+)\.+:\s(.+)$`)
var re_progress = regexp.MustCompile(`(\d+\/\d+)\s\(([\d\.]+)\%\)`)
var re_recovered = regexp.MustCompile(`(\d+\/\d+)\s\(([\d\.]+)\%\)\sDigests`)
var re_speed = regexp.MustCompile(`[\d\.]+\s\wH\/s`)
var re_mode = regexp.MustCompile(`^MODE:\s(\d+)$`)
var re_type = regexp.MustCompile(`^TYPE:\s(.+)$`)
var re_hwmon = regexp.MustCompile(`^Hardware\.Mon\.#(\d+)\.*:\s(?:(Temp):\s*(\w+))?\s?(?:(Fan):\s*(\w+)%)?\s?(?:(Util):\s*(\w+)%)?\s?(?:(Core):\s*(\w+))?\s?(?:(Mem):\s*(\w+))?\s?(?:(Bus):\s*(\w+))?$`)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

/////////

type SortByLenThenABC []string
 
func (s SortByLenThenABC) Len() int {
	return len(s)
}
 
func (s SortByLenThenABC) Less(i, j int) bool {
	len_i := len(s[i])
	len_j := len(s[j])
	return len_i < len_j || (len_i == len_j && s[i] < s[j])
}
 
func (s SortByLenThenABC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/////////

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
  b := make([]byte, length)
  for i := range b {
	b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

/////////

func StringArrayIncludes(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func intSliceToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

/////////

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

/////////

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

/////////

func spacer(w int, h int) fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{0, 0, 0, 0})
	rect.SetMinSize(fyne.Size{w, h})
	return rect
}

/*
func reserved(w int, h int) fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.RGBA{12, 5, 77, 0xff})
	rect.SetMinSize(fyne.Size{w, h})
	return rect
}
*/
