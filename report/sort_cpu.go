package report

// lessCPU compares two slaves in terms of CPU usage
func lessCPU(i *Slave, j *Slave) bool {
	return i.AllocatedResources.CPUs < j.AllocatedResources.CPUs
}
