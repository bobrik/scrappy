package report

// lessMem compares two agents in terms of memory usage
func lessMem(i *Agent, j *Agent) bool {
	return i.AllocatedResources.Memory < j.AllocatedResources.Memory
}
