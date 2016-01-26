package report

// lessTasks compares two slaves in terms of running task count
func lessTasks(i *Slave, j *Slave) bool {
	return numberOfTasks(i) < numberOfTasks(j)
}

// numberOfTasks returns number of tasks on the given slave
func numberOfTasks(slave *Slave) int {
	c := 0

	for _, r := range slave.Roles {
		c += len(r.Tasks)
	}

	return c
}
