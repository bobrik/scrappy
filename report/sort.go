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
		sort.Sort(slaveSorter{slaves: slaves, less: lessHost})
	case cpuOrder:
		sort.Sort(slaveSorter{slaves: slaves, less: lessCPU})
	case cpuPercentOrder:
		sort.Sort(slaveSorter{slaves: slaves, less: lessCPUPercent})
	case memOrder:
		sort.Sort(slaveSorter{slaves: slaves, less: lessMem})
	case memPercentOrder:
		sort.Sort(slaveSorter{slaves: slaves, less: lessMemPercent})
	case tasksOrder:
		sort.Sort(slaveSorter{slaves: slaves, less: lessTasks})
	}

	if reverse {
		for i := len(slaves)/2 - 1; i >= 0; i-- {
			opp := len(slaves) - 1 - i
			slaves[i], slaves[opp] = slaves[opp], slaves[i]
		}
	}
}

// slaveSorter is a helper structure to implement Slave sorters.
// users have to supply a slice of slaves and less fuction
// equivalent to Less() method of sort.Sorter interface
type slaveSorter struct {
	slaves []*Slave
	less   func(*Slave, *Slave) bool
}

func (s slaveSorter) Len() int {
	return len(s.slaves)
}

func (s slaveSorter) Swap(i, j int) {
	s.slaves[i], s.slaves[j] = s.slaves[j], s.slaves[i]
}

func (s slaveSorter) Less(i, j int) bool {
	return s.less(s.slaves[i], s.slaves[j])
}
