package hashcatlauncher

import (
	"os"
	"fmt"
	"strings"
	"regexp"
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
var re_checkpoint_enabled = regexp.MustCompile(`^Checkpoint enabled\. Will quit at next restore-point update\.$`)
var re_checkpoint_disabled = regexp.MustCompile(`^Checkpoint disabled\. Restore-point updates will no longer be monitored\.$`)

var re_priority = regexp.MustCompile(`^(\d+)$`)

var re_restore_file_info = regexp.MustCompile(`^hcl_(\d+)_(\d+)(?:\.restore)?$`)

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

func File_Exists(file string) bool {
	_, err := os.Stat(file);
	return err == nil
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
