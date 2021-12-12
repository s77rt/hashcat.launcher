package hashcatlauncher

import "regexp"

var reMode = regexp.MustCompile(`^Hash mode #(\d+)$`)
var reType = regexp.MustCompile(`^\s*Name\.*:\s(.+)$`)
