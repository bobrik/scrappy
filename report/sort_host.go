package report

import (
	"regexp"
	"strings"
)

// validHostname is a regexp for for hostnames like <DATACENTERC><TYPE><NUMBER>
var validHostname = regexp.MustCompile(`^(\d+)([a-z]+)(\d+)`)

// lessHost compares two slaves by comparing their hostnames
func lessHost(i *Slave, j *Slave) bool {
	if strings.Compare(i.SortString(), j.SortString()) == -1 {
		return true
	}

	return false
}
