package report

// ByTasks is a wrapper around slice of Slave instances
// to sort by number of tasks
type ByTasks struct {
	slaveSorter
}

// Less implements sort.Interface
func (s ByTasks) Less(i, j int) bool {
	return numberOfTasks(s.Slaves[i]) < numberOfTasks(s.Slaves[j])
}

// numberOfTasks returns number of tasks on the given slave
func numberOfTasks(slave *Slave) int {
	c := 0

	for _, r := range slave.Roles {
		c += len(r.Tasks)
	}

	return c
}
