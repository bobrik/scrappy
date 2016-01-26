package report

// lessMem compares two slaves in terms of memory usage
func lessMem(i *Slave, j *Slave) bool {
	return i.AllocatedResources.Memory < j.AllocatedResources.Memory
}
