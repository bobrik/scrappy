package report

import (
	"regexp"
	"strings"
)

// validHostname is a regexp for for hostnames like <DATACENTERC><TYPE><NUMBER>
var validHostname = regexp.MustCompile(`^(\d+)([a-z]+)(\d+)`)

// ByHost is a wrapper around slice of Slave instances to sort by host
type ByHost struct {
	slaveSorter
}

// Less implements sort.Interface
func (s ByHost) Less(i, j int) bool {
	if strings.Compare(s.Slaves[i].SortString(), s.Slaves[j].SortString()) == -1 {
		return true
	}

	return false
}
