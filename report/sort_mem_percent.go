package report

// lessMemPercent compares two slaves in terms of memory usage in percents
func lessMemPercent(i *Slave, j *Slave) bool {
	pi := i.AllocatedResources.Memory / i.AvailableResources.Memory
	pj := j.AllocatedResources.Memory / j.AvailableResources.Memory
	return pi < pj
}
