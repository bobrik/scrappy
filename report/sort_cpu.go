package report

// lessCPU compares two agents in terms of CPU usage
func lessCPU(i *Agent, j *Agent) bool {
	return i.AllocatedResources.CPUs < j.AllocatedResources.CPUs
}
