package report

// ByMem is a wrapper around slice of Slave instances to sort by memory usage
type ByMem struct {
	slaveSorter
}

// Less implements sort.Interface
func (s ByMem) Less(i, j int) bool {
	return s.Slaves[i].AllocatedResources.Memory < s.Slaves[j].AllocatedResources.Memory
}
