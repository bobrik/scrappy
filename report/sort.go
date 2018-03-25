package report

import "sort"

const (
	hostOrder       = "host"
	cpuOrder        = "cpu"
	cpuPercentOrder = "cpu_percent"
	memOrder        = "mem"
	memPercentOrder = "mem_percent"
	tasksOrder      = "tasks"
)

// SortAgents sorts agents in the desired order
func SortAgents(agents []*Agent, order string, reverse bool) {
	switch order {
	case hostOrder:
		sort.Sort(agentSorter{agents: agents, less: lessHost})
	case cpuOrder:
		sort.Sort(agentSorter{agents: agents, less: lessCPU})
	case cpuPercentOrder:
		sort.Sort(agentSorter{agents: agents, less: lessCPUPercent})
	case memOrder:
		sort.Sort(agentSorter{agents: agents, less: lessMem})
	case memPercentOrder:
		sort.Sort(agentSorter{agents: agents, less: lessMemPercent})
	case tasksOrder:
		sort.Sort(agentSorter{agents: agents, less: lessTasks})
	}

	if reverse {
		for i := len(agents)/2 - 1; i >= 0; i-- {
			opp := len(agents) - 1 - i
			agents[i], agents[opp] = agents[opp], agents[i]
		}
	}
}

// agentSorter is a helper structure to implement Agent sorters.
// users have to supply a slice of agents and less fuction
// equivalent to Less() method of sort.Sorter interface
type agentSorter struct {
	agents []*Agent
	less   func(*Agent, *Agent) bool
}

func (s agentSorter) Len() int {
	return len(s.agents)
}

func (s agentSorter) Swap(i, j int) {
	s.agents[i], s.agents[j] = s.agents[j], s.agents[i]
}

func (s agentSorter) Less(i, j int) bool {
	return s.less(s.agents[i], s.agents[j])
}
