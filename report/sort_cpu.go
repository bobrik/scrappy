package report

// ByCPU is a wrapper around slice of Slave instances to sort by CPU usage
type ByCPU struct {
	slaveSorter
}

// Less implements sort.Interface
func (s ByCPU) Less(i, j int) bool {
	return s.Slaves[i].AllocatedResources.CPUs < s.Slaves[j].AllocatedResources.CPUs
}
