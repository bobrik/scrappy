package report

// lessMemPercent compares two agents in terms of memory usage in percents
func lessMemPercent(i *Agent, j *Agent) bool {
	pi := i.AllocatedResources.Memory / i.AvailableResources.Memory
	pj := j.AllocatedResources.Memory / j.AvailableResources.Memory
	return pi < pj
}
