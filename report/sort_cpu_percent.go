package report

// lessCPUPercent compares two agents in terms of CPU usage in percents
func lessCPUPercent(i *Agent, j *Agent) bool {
	pi := i.AllocatedResources.CPUs / i.AvailableResources.CPUs
	pj := j.AllocatedResources.CPUs / j.AvailableResources.CPUs
	return pi < pj
}
