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

// SortSlaves sorts slaves in the desired order
func SortSlaves(slaves []*Slave, order string, reverse bool) {
	switch order {
	case hostOrder:
		sort.Sort(ByHost{slaveSorter{slaves}})
	case cpuOrder:
		sort.Sort(ByCPU{slaveSorter{slaves}})
	case cpuPercentOrder:
		sort.Sort(ByCPUPercent{slaveSorter{slaves}})
	case memOrder:
		sort.Sort(ByMem{slaveSorter{slaves}})
	case memPercentOrder:
		sort.Sort(ByMemPercent{slaveSorter{slaves}})
	case tasksOrder:
		sort.Sort(ByTasks{slaveSorter{slaves}})
	}

	if reverse {
		for i := len(slaves)/2 - 1; i >= 0; i-- {
			opp := len(slaves) - 1 - i
			slaves[i], slaves[opp] = slaves[opp], slaves[i]
		}
	}
}

// slaveSorter is a helper embeddable structure to implement Slave sorters
type slaveSorter struct {
	Slaves []*Slave
}

func (s slaveSorter) Len() int {
	return len(s.Slaves)
}

func (s slaveSorter) Swap(i, j int) {
	s.Slaves[i], s.Slaves[j] = s.Slaves[j], s.Slaves[i]
}
