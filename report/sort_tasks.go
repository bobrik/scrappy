package report

// lessTasks compares two agents in terms of running task count
func lessTasks(i *Agent, j *Agent) bool {
	return numberOfTasks(i) < numberOfTasks(j)
}

// numberOfTasks returns number of tasks on the given agent
func numberOfTasks(agent *Agent) int {
	c := 0

	for _, r := range agent.Roles {
		c += len(r.Tasks)
	}

	return c
}
