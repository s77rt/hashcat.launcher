package hashcatlauncher

import (
	"os"
	"fmt"
	"strconv"
	"strings"
	"path/filepath"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/s77rt/hashcat.launcher/pkg/subprocess"
	"github.com/s77rt/hashcat.launcher/pkg/xfyne/xwidget"
)

// Types
type hashcat struct {
	binary_file string
	available_hash_types []*xwidget.SelectorOption
	args hashcat_args
}

type hashcat_args struct {
	session string
	attack_mode hashcat_attack_mode
	hash_file string
	hash_type hashcat_hash_type
	separator string
	remove_found_hashes bool
	disable_potfile bool
	ignore_usernames bool
	disable_monitor bool
	temp_abort int
	devices_types []int
	workload_profile hashcat_workload_profile
	outfile string
	outfile_format []int
	optimized_kernel bool
	slower_candidate bool
	force bool
	status_timer int
	attack_payload string
}

type hashcat_workload_profile int
const (
	_ hashcat_workload_profile = iota
	hashcat_workload_profile_Low
	hashcat_workload_profile_Default
	hashcat_workload_profile_High
	hashcat_workload_profile_Nightmare
)

type hashcat_attack_mode int
const (
	hashcat_attack_mode_Dictionary hashcat_attack_mode = iota
	hashcat_attack_mode_Combinator
	_
	hashcat_attack_mode_Mask
	_
	_
	hashcat_attack_mode_Hybrid1
	hashcat_attack_mode_Hybrid2
)

type hashcat_hash_type int

// Get Functions
func get_available_hash_typess(hcl_gui *hcl_gui) {
	var hashmode, hashtype string
	hcl_gui.hashcat.available_hash_types = []*xwidget.SelectorOption{}
	wdir, _ := filepath.Split(hcl_gui.hashcat.binary_file)
	cmd := subprocess.Subprocess{
		subprocess.SubprocessStatusNotRunning,
		wdir,
		hcl_gui.hashcat.binary_file,
		[]string{"--example-hashes", "--quiet"},
		nil,
		nil,
		func(s string) {
			mode_line := re_mode.FindStringSubmatch(s)
			if len(mode_line) == 2 {
				hashmode = mode_line[1]
			} else {
				type_line := re_type.FindStringSubmatch(s)
				if len(type_line) == 2 {
					hashtype = type_line[1]
					hcl_gui.hashcat.available_hash_types = append(hcl_gui.hashcat.available_hash_types, xwidget.NewSelectorOptionWithStyle(fmt.Sprintf("%s - %s", hashmode, hashtype), hashmode, fyne.TextAlignLeading, fyne.TextStyle{}, func(string){}))
				}
			}
		},
		func(s string) {
			fmt.Fprintf(os.Stderr, "%s\n", s)
		},
		func(){},
		func(){},
	}
	cmd.Execute()
}

func get_devices_info(hcl_gui *hcl_gui) string {
	info := ""
	errors := ""
	wdir, _ := filepath.Split(hcl_gui.hashcat.binary_file)
	cmd := subprocess.Subprocess{
		subprocess.SubprocessStatusNotRunning,
		wdir,
		hcl_gui.hashcat.binary_file,
		[]string{"-I", "--force", "--quiet"},
		nil,
		nil,
		func(s string) {
			info += re_ansi.ReplaceAllString(s, "")+"\n"
		},
		func(s string) {
			fmt.Fprintf(os.Stderr, "%s\n", s)
			errors += re_ansi.ReplaceAllString(s, "")+"\n"
		},
		func(){},
		func(){},
	}
	cmd.Execute()
	if len(errors) > 0 {
		info += "\nErrors:\n"+errors
	}
	return info
}

func get_benchmark(hcl_gui *hcl_gui) string {
	benchmark := ""
	errors := ""
	wdir, _ := filepath.Split(hcl_gui.hashcat.binary_file)
	cmd := subprocess.Subprocess{
		subprocess.SubprocessStatusNotRunning,
		wdir,
		hcl_gui.hashcat.binary_file,
		[]string{"-O", fmt.Sprintf("-m%d", hcl_gui.hashcat.args.hash_type), "-b", fmt.Sprintf("-w%d", hcl_gui.hashcat.args.workload_profile), "-D", intSliceToString(hcl_gui.hashcat.args.devices_types,","), "--force", "--quiet"},
		nil,
		nil,
		func(s string) {
			benchmark += re_ansi.ReplaceAllString(s, "")+"\n"
		},
		func(s string) {
			fmt.Fprintf(os.Stderr, "%s\n", s)
			errors += re_ansi.ReplaceAllString(s, "")+"\n"
		},
		func(){},
		func(){},
	}
	cmd.Execute()
	if len(errors) > 0 {
		benchmark += "\nErrors:\n"+errors
	}
	return benchmark
}

// Set Functions
func set_hash_type(hcl_gui *hcl_gui, hash_type_fakeselector *widget.Box, value string) {
	value_int, _ := strconv.ParseInt(value, 10, 32)
	hcl_gui.hashcat.args.hash_type = hashcat_hash_type(value_int)
	fake_hash_type_selector_hack(hcl_gui, hash_type_fakeselector, value)
}

func set_attack_mode(hcl_gui *hcl_gui, value string) {
	switch value{
	case "Dictionary":
		hcl_gui.hc_attack_mode.Selected = "Dictionary"
		hcl_gui.hashcat.args.attack_mode = hashcat_attack_mode_Dictionary
		hcl_gui.hc_dictionary_attack_conf.Show()
		hcl_gui.hc_combinator_attack_conf.Hide()
		hcl_gui.hc_mask_attack_conf.Hide()
		hcl_gui.hc_hybrid1_attack_conf.Hide()
		hcl_gui.hc_hybrid2_attack_conf.Hide()
	case "Combinator":
		hcl_gui.hc_attack_mode.Selected = "Combinator"
		hcl_gui.hashcat.args.attack_mode = hashcat_attack_mode_Combinator
		hcl_gui.hc_dictionary_attack_conf.Hide()
		hcl_gui.hc_combinator_attack_conf.Show()
		hcl_gui.hc_mask_attack_conf.Hide()
		hcl_gui.hc_hybrid1_attack_conf.Hide()
		hcl_gui.hc_hybrid2_attack_conf.Hide()
	case "Mask":
		hcl_gui.hc_attack_mode.Selected = "Mask"
		hcl_gui.hashcat.args.attack_mode = hashcat_attack_mode_Mask
		hcl_gui.hc_dictionary_attack_conf.Hide()
		hcl_gui.hc_combinator_attack_conf.Hide()
		hcl_gui.hc_mask_attack_conf.Show()
		hcl_gui.hc_hybrid1_attack_conf.Hide()
		hcl_gui.hc_hybrid2_attack_conf.Hide()
	case "Hybrid1 (Dict+Mask)":
		hcl_gui.hc_attack_mode.Selected = "Hybrid1 (Dict+Mask)"
		hcl_gui.hashcat.args.attack_mode = hashcat_attack_mode_Hybrid1
		hcl_gui.hc_dictionary_attack_conf.Hide()
		hcl_gui.hc_combinator_attack_conf.Hide()
		hcl_gui.hc_mask_attack_conf.Hide()
		hcl_gui.hc_hybrid1_attack_conf.Show()
		hcl_gui.hc_hybrid2_attack_conf.Hide()
	case "Hybrid2 (Mask+Dict)":
		hcl_gui.hc_attack_mode.Selected = "Hybrid2 (Mask+Dict)"
		hcl_gui.hashcat.args.attack_mode = hashcat_attack_mode_Hybrid2
		hcl_gui.hc_dictionary_attack_conf.Hide()
		hcl_gui.hc_combinator_attack_conf.Hide()
		hcl_gui.hc_mask_attack_conf.Hide()
		hcl_gui.hc_hybrid1_attack_conf.Hide()
		hcl_gui.hc_hybrid2_attack_conf.Show()
	}
}

func set_hash_file(hcl_gui *hcl_gui, file string) {
	hcl_gui.hashcat.args.hash_file = file
}

func set_separator(hcl_gui *hcl_gui) {
	hcl_gui.hashcat.args.separator = string(hcl_gui.hc_separator.Text[3])
}

func set_remove_found_hashes(hcl_gui *hcl_gui, check bool) {
	hcl_gui.hashcat.args.remove_found_hashes = check
}

func set_disable_potfile(hcl_gui *hcl_gui, check bool) {
	hcl_gui.hashcat.args.disable_potfile = check
}

func set_ignore_usernames(hcl_gui *hcl_gui, check bool) {
	hcl_gui.hashcat.args.ignore_usernames = check
}

func set_disable_monitor(hcl_gui *hcl_gui, check bool) {
	if check {
		hcl_gui.hc_temp_abort.Selected = "OFF"
		hcl_gui.hc_temp_abort.Refresh()
	} else {
		hcl_gui.hc_temp_abort.SetSelected("90")
	}
	hcl_gui.hashcat.args.disable_monitor = check
}

func set_temp_abort(hcl_gui *hcl_gui, temp string) {
	hcl_gui.hashcat.args.disable_monitor = false
	hcl_gui.hc_disable_monitor.Checked = false
	hcl_gui.hc_disable_monitor.Refresh()
	temp_int, _ := strconv.ParseInt(temp, 10, 32)
	hcl_gui.hashcat.args.temp_abort = int(temp_int)
}

func set_devices_types(hcl_gui *hcl_gui, devices string) {
	var devices_int []int
	switch devices {
	case "GPU":
		devices_int = []int{2}
	case "CPU":
		devices_int = []int{1}
	case "FPGA":
		devices_int = []int{3}
	case "GPU+CPU":
		devices_int = []int{1,2}
	case "GPU+FPGA":
		devices_int = []int{2,3}
	case "CPU+FPGA":
		devices_int = []int{1,3}
	case "All":
		devices_int = []int{1,2,3}
	}
	hcl_gui.hashcat.args.devices_types = devices_int
}

func set_workload_profile(hcl_gui *hcl_gui, profile string) {
	var workload_profile hashcat_workload_profile
	switch profile {
	case "Low":
		workload_profile = hashcat_workload_profile_Low
	case "Default":
		workload_profile = hashcat_workload_profile_Default
	case "High":
		workload_profile = hashcat_workload_profile_High
	case "Nightmare":
		workload_profile = hashcat_workload_profile_Nightmare
	}
	hcl_gui.hashcat.args.workload_profile = workload_profile
}

func set_outfile(hcl_gui *hcl_gui, outfile string) {
	hcl_gui.hashcat.args.outfile = outfile
}

func set_optimized_kernel(hcl_gui *hcl_gui, check bool) {
	hcl_gui.hashcat.args.optimized_kernel = check
}

func set_slower_candidate(hcl_gui *hcl_gui, check bool) {
	hcl_gui.hashcat.args.slower_candidate = check
}

func set_force(hcl_gui *hcl_gui, check bool) {
	hcl_gui.hashcat.args.force = check
}

func set_outfile_format(hcl_gui *hcl_gui, outfile_format string) {
	var outfile_format_int []int
	switch outfile_format {
	case "hash[:salt]:plain":
		outfile_format_int = []int{1,2}
	case "hash[:salt]:plain:hex_plain":
		outfile_format_int = []int{1,2,3}
	/*
	case "hash[:salt]:plain:crack_pos":
		outfile_format_int = []int{1,2,4}
	case "hash[:salt]:plain:hex_plain:crack_pos":
		outfile_format_int = []int{1,2,3,4}
	*/
	}
	hcl_gui.hashcat.args.outfile_format = outfile_format_int
}

// Others
func get_mask_length(mask string) int {
	if len(mask) == 0 {
		return 0
	}
	length := 0
	skip_next := false
	for _, s := range strings.Split(mask, " ") {
		if len(s) == 0 {
			continue
		}
		if skip_next == true {
			skip_next = false
			continue
		}
		if s[0] == 0x2d {
			if len(s) == 2 {
				skip_next = true
			}
			continue
		}
		for _, l := range s {
			if l == 0x3f {
				continue
			}
			length++
		}
		if length > 0 {
			break
		}
	}
	return length
}

// Init
func hashcat_init(hcl_gui *hcl_gui) {
	hcl_gui.hc_attack_mode.Options = []string{"Dictionary", "Combinator", "Mask", "Hybrid1 (Dict+Mask)", "Hybrid2 (Mask+Dict)"}
	hcl_gui.hc_attack_mode.SetSelected("Dictionary")

	hcl_gui.hashcat.args.hash_type = -1
	go get_available_hash_typess(hcl_gui)

	hcl_gui.hc_temp_abort.SetSelected("90")

	hcl_gui.hc_devices_types.SetSelected("GPU")

	hcl_gui.hc_wordload_profiles.SetSelected("Default")
}
