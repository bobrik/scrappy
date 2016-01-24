package report

// ByMemPercent is a wrapper around slice of Slave instances
// to sort by memory usage in percents
type ByMemPercent struct {
	slaveSorter
}

// Less implements sort.Interface
func (s ByMemPercent) Less(i, j int) bool {
	pi := s.Slaves[i].AllocatedResources.Memory / s.Slaves[i].AvailableResources.Memory
	pj := s.Slaves[j].AllocatedResources.Memory / s.Slaves[j].AvailableResources.Memory
	return pi < pj
}
