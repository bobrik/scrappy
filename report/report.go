package report

import (
	"fmt"
	"log"

	"github.com/bobrik/scrappy/mesos"
)

// Report contains resource usage on slave, role and task leve
type Report struct {
	Slaves []*Slave `json:"slaves"`
}

// Slave represents slave's state in the report
type Slave struct {
	ID                 string                 `json:"id"`
	Hostname           string                 `json:"hostname"`
	Attributes         map[string]interface{} `json:"attributes"`
	AvailableResources mesos.Resources        `json:"available_resources"`
	AllocatedResources mesos.Resources        `json:"allocated_resources"`
	Roles              map[string]*Role       `json:"roles"`
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

// Generate converts Mesos state into the report
func Generate(state *mesos.State, role string) *Report {
	r := &Report{
		Slaves: make([]*Slave, 0, len(state.Slaves)),
	}

	slaves := make(map[string]*Slave, len(state.Slaves))

	for _, slave := range slaveMap(state) {
		slaves[slave.ID] = &Slave{
			ID:                 slave.ID,
			Hostname:           slave.Hostname,
			Attributes:         slave.Attributes,
			AvailableResources: slave.Resources,
			Roles:              make(map[string]*Role, len(slave.ReservedResources)+1),
		}

		if slave.UnreservedResources.CPUs > 0 && slave.UnreservedResources.Memory > 0 {
			slaves[slave.ID].Roles["*"] = &Role{
				Name:               "*",
				Tasks:              []*Task{},
				AvailableResources: slave.UnreservedResources,
			}
		}

		for name, resources := range slave.ReservedResources {
			slaves[slave.ID].Roles[name] = &Role{
				Name:               name,
				Tasks:              []*Task{},
				AvailableResources: resources,
			}
		}
	}

	for _, f := range state.Frameworks {
		for _, t := range f.Tasks {
			slave := slaves[t.SlaveID]
			role := f.Role
			if t.Role != "" {
				role = t.Role
			}
			// marathon grew multi_role support, and is subscribed implicitly to '*' in addition
			// to other roles it may subscribe from.
			// However, the 'role' that is set at the task level isn't accurate in the case where
			// the task is actually using '*' resources.  Log a warning that stats may be inaccurate
			// and add it to '*' if this occurs.  Longer term once mesos has fixed this, this should
			// be removed.
			if _, ok := slave.Roles[role]; !ok {
				log.Printf("warning: results will be inaccurate: task id %s on slave id %s (host %s) has role %s which isn't a defined resource; assuming it's an implicit '*' instead.", t.ID, t.SlaveID, slave.Hostname, role)
				role = "*"
			}

			slave.AllocatedResources.Add(t.Resources)
			slave.Roles[role].AllocatedResources.Add(t.Resources)

			slave.Roles[role].Tasks = append(slave.Roles[role].Tasks, &Task{
				ID:        t.ID,
				Name:      t.Name,
				Framework: f.Name,
				Resources: t.Resources,
			})
		}
	}

	filter(slaves, role)

	for _, s := range slaves {
		r.Slaves = append(r.Slaves, s)
	}

	return r
}

// filter removes slaves that do not include given role
func filter(slaves map[string]*Slave, role string) {
	if role == "" {
		return
	}

	for id, slave := range slaves {
		if _, ok := slave.Roles[role]; !ok {
			delete(slaves, id)
		}
	}
}

// slaveMap converts a slice of slaves into a map id -> slave
func slaveMap(state *mesos.State) map[string]mesos.Slave {
	m := make(map[string]mesos.Slave, len(state.Slaves))

	for _, s := range state.Slaves {
		m[s.ID] = s
	}

	return m
}
