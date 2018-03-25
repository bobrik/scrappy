package report

import (
	"fmt"

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

// Attribute represents a key-value attribute in the report
type Attribute struct {
	Key                string          `json:"key"`
	Value              interface{}     `json:"value"`
	AvailableResources mesos.Resources `json:"available_resources"`
	AllocatedResources mesos.Resources `json:"allocated_resources"`
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
