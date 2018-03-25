package report

import (
	"fmt"

	"github.com/bobrik/scrappy/mesos"
)

// Report contains resource usage on agent, role and task leve
type Report struct {
	Agents []*Agent `json:"slaves"`
}

// Agent represents agent's state in the report
type Agent struct {
	ID                 string                 `json:"id"`
	Hostname           string                 `json:"hostname"`
	Attributes         map[string]interface{} `json:"attributes"`
	AvailableResources mesos.Resources        `json:"available_resources"`
	AllocatedResources mesos.Resources        `json:"allocated_resources"`
	Roles              map[string]*Role       `json:"roles"`
}

// SortString returns string representation for sorting
func (s Agent) SortString() string {
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
		Agents: make([]*Agent, 0, len(state.Agents)),
	}

	agents := make(map[string]*Agent, len(state.Agents))

	for _, agent := range agentMap(state) {
		agents[agent.ID] = &Agent{
			ID:                 agent.ID,
			Hostname:           agent.Hostname,
			Attributes:         agent.Attributes,
			AvailableResources: agent.Resources,
			Roles:              make(map[string]*Role, len(agent.ReservedResources)+1),
		}

		if agent.UnreservedResources.CPUs > 0 && agent.UnreservedResources.Memory > 0 {
			agents[agent.ID].Roles["*"] = &Role{
				Name:               "*",
				Tasks:              []*Task{},
				AvailableResources: agent.UnreservedResources,
			}
		}

		for name, resources := range agent.ReservedResources {
			agents[agent.ID].Roles[name] = &Role{
				Name:               name,
				Tasks:              []*Task{},
				AvailableResources: resources,
			}
		}
	}

	for _, f := range state.Frameworks {
		for _, t := range f.Tasks {
			agent := agents[t.AgentID]

			agent.AllocatedResources.Add(t.Resources)
			agent.Roles[f.Role].AllocatedResources.Add(t.Resources)

			agent.Roles[f.Role].Tasks = append(agent.Roles[f.Role].Tasks, &Task{
				ID:        t.ID,
				Name:      t.Name,
				Framework: f.Name,
				Resources: t.Resources,
			})
		}
	}

	filter(agents, role)

	for _, s := range agents {
		r.Agents = append(r.Agents, s)
	}

	return r
}

// filter removes agents that do not include given role
func filter(agents map[string]*Agent, role string) {
	if role == "" {
		return
	}

	for id, agent := range agents {
		if _, ok := agent.Roles[role]; !ok {
			delete(agents, id)
		}
	}
}

// agentMap converts a slice of agents into a map id -> agent
func agentMap(state *mesos.State) map[string]mesos.Agent {
	m := make(map[string]mesos.Agent, len(state.Agents))

	for _, s := range state.Agents {
		m[s.ID] = s
	}

	return m
}
