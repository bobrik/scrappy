package report

// ByCPUPercent is a wrapper around slice of Slave instances
// to sort by CPU usage in percents
type ByCPUPercent struct {
	slaveSorter
}

// Less implements sort.Interface
func (s ByCPUPercent) Less(i, j int) bool {
	pi := s.Slaves[i].AllocatedResources.CPUs / s.Slaves[i].AvailableResources.CPUs
	pj := s.Slaves[j].AllocatedResources.CPUs / s.Slaves[j].AvailableResources.CPUs
	return pi < pj
}
