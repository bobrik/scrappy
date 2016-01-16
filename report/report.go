package report

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/bobrik/scrappy/mesos"
)

// validHostname is a regexp for for hostnames like <DATACENTERC><TYPE><NUMBER>
var validHostname = regexp.MustCompile(`^(\d+)([a-z]+)(\d+)`)

// Report contains resource usage on slave, role and task leve
type Report struct {
	Slaves Slaves `json:"slaves"`
}

// Slave represents slave's state in the report
type Slave struct {
	ID                 string           `json:"id"`
	Hostname           string           `json:"hostname"`
	AvailableResources mesos.Resources  `json:"available_resources"`
	AllocatedResources mesos.Resources  `json:"allocated_resources"`
	Roles              map[string]*Role `json:"roles"`
}

// SortString returns string representation for sorting
func (s Slave) SortString() string {
	p := validHostname.FindStringSubmatch(s.Hostname)

	if p == nil {
		return ""
	}

	return fmt.Sprintf("%010s%s%010s", p[1], p[2], p[3])
}

// Role represents role's state in the report
type Role struct {
	Name               string          `json:"name"`
	Tasks              []*Task         `json:"tasks"`
	AvailableResources mesos.Resources `json:"available_resources"`
	AllocatedResources mesos.Resources `json:"allocated_resources"`
}

// Task represents task's state in the report
type Task struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Framework string          `json:"string"`
	Resources mesos.Resources `json:"resources"`
}

// Slaves is a slice of Slave instances
type Slaves []*Slave

func (s Slaves) Len() int {
	return len(s)
}

func (s Slaves) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Slaves) Less(i, j int) bool {
	if strings.Compare(s[i].SortString(), s[j].SortString()) == -1 {
		return true
	}

	return false
}

// Generate converts Mesos state into the report
func Generate(state *mesos.State) *Report {
	r := &Report{
		Slaves: make([]*Slave, 0, len(state.Slaves)),
	}

	sm := make(map[string]*Slave, len(state.Slaves))

	for _, slave := range slaveMap(state) {
		sm[slave.ID] = &Slave{
			ID:                 slave.ID,
			Hostname:           slave.Hostname,
			AvailableResources: slave.Resources,
			Roles:              make(map[string]*Role, len(slave.ReservedResources)+1),
		}

		if slave.UnreservedResources.CPUs > 0 && slave.UnreservedResources.Memory > 0 {
			sm[slave.ID].Roles["*"] = &Role{
				Name:               "*",
				Tasks:              []*Task{},
				AvailableResources: slave.UnreservedResources,
			}
		}

		for name, resources := range slave.ReservedResources {
			sm[slave.ID].Roles[name] = &Role{
				Name:               name,
				Tasks:              []*Task{},
				AvailableResources: resources,
			}
		}
	}

	for _, f := range state.Frameworks {
		for _, t := range f.Tasks {
			slave := sm[t.SlaveID]

			slave.AllocatedResources.Add(t.Resources)
			slave.Roles[f.Role].AllocatedResources.Add(t.Resources)

			slave.Roles[f.Role].Tasks = append(slave.Roles[f.Role].Tasks, &Task{
				ID:        t.ID,
				Name:      t.Name,
				Framework: f.Name,
				Resources: t.Resources,
			})
		}
	}

	for _, s := range sm {
		r.Slaves = append(r.Slaves, s)
	}

	sort.Sort(r.Slaves)

	return r
}

// slaveMap converts a slice of slaves into a map id -> slave
func slaveMap(state *mesos.State) map[string]mesos.Slave {
	m := make(map[string]mesos.Slave, len(state.Slaves))

	for _, s := range state.Slaves {
		m[s.ID] = s
	}

	return m
}
