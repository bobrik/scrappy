package report

import (
	"regexp"
	"strings"
)

// validHostname is a regexp for for hostnames like <DATACENTERC><TYPE><NUMBER>
var validHostname = regexp.MustCompile(`^(\d+)([a-z]+)(\d+)`)

// lessHost compares two agents by comparing their hostnames
func lessHost(i *Agent, j *Agent) bool {
	if strings.Compare(i.SortString(), j.SortString()) == -1 {
		return true
	}

	return false
}
