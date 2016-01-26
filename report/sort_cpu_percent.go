package report

// lessCPUPercent compares two slaves in terms of CPU usage in percents
func lessCPUPercent(i *Slave, j *Slave) bool {
	pi := i.AllocatedResources.CPUs / i.AvailableResources.CPUs
	pj := j.AllocatedResources.CPUs / j.AvailableResources.CPUs
	return pi < pj
}
