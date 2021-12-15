package hashcatlauncher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/s77rt/hashcat.launcher/pkg/subprocess"
)

type Hashcat struct {
	BinaryFile string           `json:"binaryFile"`
	Algorithms map[int64]string `json:"algorithms"`
}

var DefaultSessionID = "hashcat"

func (h *Hashcat) GetAlgorithms() {
	h.Algorithms = make(map[int64]string)

	var algorithmMode int64
	var algorithmName string

	args := []string{"--hash-info", "--quiet"}
	wdir, _ := filepath.Split(h.BinaryFile)
	cmd := subprocess.Subprocess{
		subprocess.SubprocessStatusNotRunning,
		wdir,
		h.BinaryFile,
		args,
		nil,
		nil,
		func(s string) {
			modeLine := reMode.FindStringSubmatch(s)
			if len(modeLine) == 2 {
				var err error
				algorithmMode, err = strconv.ParseInt(modeLine[1], 10, 64)
				if err != nil {
					return
				}
			} else {
				typeLine := reType.FindStringSubmatch(s)
				if len(typeLine) == 2 {
					algorithmName = typeLine[1]
					h.Algorithms[algorithmMode] = algorithmName
				}
			}
		},
		func(s string) {
			fmt.Fprintf(os.Stderr, "%s\n", s)
		},
		func() {},
		func() {},
	}
	cmd.Execute()
}

type HashcatArgs struct {
	Session *string

	AttackMode *HashcatAttackMode `json:"attackMode"`
	HashMode   *HashcatHashMode   `json:"hashMode"`

	Dictionaries    *[]string `json:"dictionaries"`    // Files
	Rules           *[]string `json:"rules"`           // Files
	Mask            *string   `json:"mask"`            // Direct Input
	MaskFile        *string   `json:"maskFile"`        // File
	LeftDictionary  *string   `json:"leftDictionary"`  // File
	LeftRule        *string   `json:"leftRule"`        // Direct Input
	RightDictionary *string   `json:"rightDictionary"` // File
	RightRule       *string   `json:"rightRule"`       // Direct Input

	CustomCharset1 *string `json:"customCharset1"`
	CustomCharset2 *string `json:"customCharset2"`
	CustomCharset3 *string `json:"customCharset3"`
	CustomCharset4 *string `json:"customCharset4"`

	EnableMaskIncrementMode *bool  `json:"enableMaskIncrementMode"`
	MaskIncrementMin        *int64 `json:"maskIncrementMin"`
	MaskIncrementMax        *int64 `json:"maskIncrementMax"`

	Hash *string `json:"hash"` // File

	EnableOptimizedKernel           *bool `json:"enableOptimizedKernel"`
	EnableSlowerCandidateGenerators *bool `json:"enableSlowerCandidateGenerators"`
	RemoveFoundHashes               *bool `json:"removeFoundHashes"`
	DisablePotFile                  *bool `json:"disablePotFile"`
	IgnoreUsernames                 *bool `json:"ignoreUsernames"`
	DisableSelfTest                 *bool `json:"disableSelfTest"`
	IgnoreWarnings                  *bool `json:"ignoreWarnings"`

	DevicesIDs      *[]int64 `json:"devicesIDs"`
	DevicesTypes    *[]int64 `json:"devicesTypes"`
	WorkloadProfile *int64   `json:"workloadProfile"`

	DisableMonitor *bool  `json:"disableMonitor"`
	TempAbort      *int64 `json:"tempAbort"`

	MarkovDisable   *bool  `json:"markovDisable"`
	MarkovClassic   *bool  `json:"markovClassic"`
	MarkovThreshold *int64 `json:"markovThreshold"`

	ExtraArguments *[]string `json:"extraArguments"`

	StatusTimer *int64 `json:"statusTimer"`

	OutputFile   *string  `json:"outputFile"`
	OutputFormat *[]int64 `json:"outputFormat"`
}

func (ha *HashcatArgs) Build() (args []string, err error) {
	if ha.Session == nil {
		ha.Session = &DefaultSessionID
	}

	if ha.Hash == nil {
		err = errors.New("Missing hash")
		return
	}

	if ha.HashMode == nil {
		err = errors.New("Missing hash mode (algorithm)")
		return
	}

	if ha.AttackMode == nil {
		err = errors.New("Missing attack mode")
		return
	}

	if ha.StatusTimer == nil {
		err = errors.New("Missing status timer")
		return
	}

	if ha.OutputFile == nil {
		err = errors.New("Missing output file")
		return
	}

	if ha.OutputFormat == nil {
		err = errors.New("Missing output format")
		return
	}

	args = append(args, fmt.Sprintf("--session=%s", *ha.Session))

	args = append(args, []string{"--status", "--status-json", fmt.Sprintf("--status-timer=%d", *ha.StatusTimer)}...)

	if ha.EnableOptimizedKernel != nil && *ha.EnableOptimizedKernel == true {
		args = append(args, "-O")
	}

	if ha.EnableSlowerCandidateGenerators != nil && *ha.EnableSlowerCandidateGenerators == true {
		args = append(args, "-S")
	}

	if ha.RemoveFoundHashes != nil && *ha.RemoveFoundHashes == true {
		args = append(args, "--remove")
	}

	if ha.DisablePotFile != nil && *ha.DisablePotFile == true {
		args = append(args, "--potfile-disable")
	}

	if ha.IgnoreUsernames != nil && *ha.IgnoreUsernames == true {
		args = append(args, "--username")
	}

	if ha.DisableSelfTest != nil && *ha.DisableSelfTest == true {
		args = append(args, "--self-test-disable")
	}

	if ha.IgnoreWarnings != nil && *ha.IgnoreWarnings == true {
		args = append(args, "--force")
	}

	if ha.DisableMonitor != nil && *ha.DisableMonitor == true {
		args = append(args, "--hwmon-disable")
	} else if ha.TempAbort != nil {
		args = append(args, fmt.Sprintf("--hwmon-temp-abort=%d", *ha.TempAbort))
	}

	if ha.MarkovDisable != nil && *ha.MarkovDisable == true {
		args = append(args, "--markov-disable")
	}
	if ha.MarkovClassic != nil && *ha.MarkovClassic == true {
		args = append(args, "--markov-classic")
	}
	if ha.MarkovThreshold != nil {
		args = append(args, fmt.Sprintf("--markov-threshold=%d", *ha.MarkovThreshold))
	}

	if ha.WorkloadProfile != nil {
		args = append(args, fmt.Sprintf("-w%d", *ha.WorkloadProfile))
	}

	args = append(args, fmt.Sprintf("-m%d", *ha.HashMode))
	args = append(args, fmt.Sprintf("-a%d", *ha.AttackMode))
	args = append(args, *ha.Hash)

	if ha.DevicesIDs != nil {
		args = append(args, []string{"-d", strings.Trim(strings.Replace(fmt.Sprint(*ha.DevicesIDs), " ", ",", -1), "[]")}...)
	}

	if ha.DevicesTypes != nil {
		args = append(args, []string{"-D", strings.Trim(strings.Replace(fmt.Sprint(*ha.DevicesTypes), " ", ",", -1), "[]")}...)
	}

	if ha.ExtraArguments != nil && len(*ha.ExtraArguments) > 0 {
		args = append(args, *ha.ExtraArguments...)
	}

	args = append(args, []string{"-o", *ha.OutputFile}...)
	args = append(args, fmt.Sprintf("--outfile-format=%s", strings.Trim(strings.Replace(fmt.Sprint(*ha.OutputFormat), " ", ",", -1), "[]")))

	switch *ha.AttackMode {
	case HashcatAttackModeDictionary:
		if ha.Dictionaries == nil {
			err = errors.New("Missing dictionaries")
			return
		}
		args = append(args, *ha.Dictionaries...)
		if ha.Rules != nil {
			for _, rule := range *ha.Rules {
				args = append(args, []string{"-r", rule}...)
			}
		}
	case HashcatAttackModeCombinator:
		if ha.LeftDictionary == nil {
			err = errors.New("Missing left dictionary")
			return
		}
		if ha.RightDictionary == nil {
			err = errors.New("Missing right dictionary")
			return
		}
		args = append(args, []string{*ha.LeftDictionary, *ha.RightDictionary}...)
		if ha.LeftRule != nil {
			args = append(args, []string{"-j", *ha.LeftRule}...)
		}
		if ha.RightRule != nil {
			args = append(args, []string{"-k", *ha.RightRule}...)
		}
	case HashcatAttackModeMask:
		if ha.MaskFile != nil {
			args = append(args, *ha.MaskFile)
		} else if ha.Mask != nil {
			if ha.CustomCharset1 != nil {
				args = append(args, []string{"-1", *ha.CustomCharset1}...)
			}
			if ha.CustomCharset2 != nil {
				args = append(args, []string{"-2", *ha.CustomCharset2}...)
			}
			if ha.CustomCharset3 != nil {
				args = append(args, []string{"-3", *ha.CustomCharset3}...)
			}
			if ha.CustomCharset4 != nil {
				args = append(args, []string{"-4", *ha.CustomCharset4}...)
			}
			args = append(args, *ha.Mask)
		} else {
			err = errors.New("Missing mask")
			return
		}
		if ha.EnableMaskIncrementMode != nil && *ha.EnableMaskIncrementMode == true {
			if ha.MaskIncrementMin == nil || ha.MaskIncrementMax == nil {
				err = errors.New("Missing mask increment min and/or max")
				return
			}
			if *ha.MaskIncrementMin > *ha.MaskIncrementMax {
				err = errors.New("mask increment min cannot be greater than mask increment max")
				return
			}
			args = append(args, []string{"-i", fmt.Sprintf("--increment-min=%d", *ha.MaskIncrementMin), fmt.Sprintf("--increment-max=%d", *ha.MaskIncrementMax)}...)
		}
	case HashcatAttackModeHybrid1:
		// Left (Dictionary)
		if ha.LeftDictionary == nil {
			err = errors.New("Missing dictionary")
			return
		}
		args = append(args, *ha.LeftDictionary)
		if ha.LeftRule != nil {
			args = append(args, []string{"-j", *ha.LeftRule}...)
		}
		// Right (Mask)
		if ha.MaskFile != nil {
			args = append(args, *ha.MaskFile)
		} else if ha.Mask != nil {
			if ha.CustomCharset1 != nil {
				args = append(args, []string{"-1", *ha.CustomCharset1}...)
			}
			if ha.CustomCharset2 != nil {
				args = append(args, []string{"-2", *ha.CustomCharset2}...)
			}
			if ha.CustomCharset3 != nil {
				args = append(args, []string{"-3", *ha.CustomCharset3}...)
			}
			if ha.CustomCharset4 != nil {
				args = append(args, []string{"-4", *ha.CustomCharset4}...)
			}
			args = append(args, *ha.Mask)
		} else {
			err = errors.New("Missing mask")
			return
		}
		if ha.EnableMaskIncrementMode != nil && *ha.EnableMaskIncrementMode == true {
			if ha.MaskIncrementMin == nil || ha.MaskIncrementMax == nil {
				err = errors.New("Missing mask increment min and/or max")
				return
			}
			if *ha.MaskIncrementMin > *ha.MaskIncrementMax {
				err = errors.New("mask increment min cannot be greater than mask increment max")
				return
			}
			args = append(args, []string{"-i", fmt.Sprintf("--increment-min=%d", *ha.MaskIncrementMin), fmt.Sprintf("--increment-max=%d", *ha.MaskIncrementMax)}...)
		}
	case HashcatAttackModeHybrid2:
		// Left (Mask)
		if ha.MaskFile != nil {
			args = append(args, *ha.MaskFile)
		} else if ha.Mask != nil {
			if ha.CustomCharset1 != nil {
				args = append(args, []string{"-1", *ha.CustomCharset1}...)
			}
			if ha.CustomCharset2 != nil {
				args = append(args, []string{"-2", *ha.CustomCharset2}...)
			}
			if ha.CustomCharset3 != nil {
				args = append(args, []string{"-3", *ha.CustomCharset3}...)
			}
			if ha.CustomCharset4 != nil {
				args = append(args, []string{"-4", *ha.CustomCharset4}...)
			}
			args = append(args, *ha.Mask)
		} else {
			err = errors.New("Missing mask")
			return
		}
		if ha.EnableMaskIncrementMode != nil && *ha.EnableMaskIncrementMode == true {
			if ha.MaskIncrementMin == nil || ha.MaskIncrementMax == nil {
				err = errors.New("Missing mask increment min and/or max")
				return
			}
			if *ha.MaskIncrementMin > *ha.MaskIncrementMax {
				err = errors.New("mask increment min cannot be greater than mask increment max")
				return
			}
			args = append(args, []string{"-i", fmt.Sprintf("--increment-min=%d", *ha.MaskIncrementMin), fmt.Sprintf("--increment-max=%d", *ha.MaskIncrementMax)}...)
		}
		// Right (Dictionary)
		if ha.RightDictionary == nil {
			err = errors.New("Missing dictionary")
			return
		}
		args = append(args, *ha.RightDictionary)
		if ha.RightRule != nil {
			args = append(args, []string{"-k", *ha.RightRule}...)
		}
	default:
		err = errors.New("Unsupported attack mode")
		return
	}

	return
}

type HashcatAttackMode int64

const (
	HashcatAttackModeDictionary HashcatAttackMode = iota
	HashcatAttackModeCombinator
	_
	HashcatAttackModeMask
	_
	_
	HashcatAttackModeHybrid1
	HashcatAttackModeHybrid2
)

type HashcatHashMode int64
