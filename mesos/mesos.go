package mesos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// State represents Mesos cluster state
type State struct {
	Frameworks []Framework `json:"frameworks"`
	Slaves     []Slave     `json:"slaves"`
}

// Framework represents framework info from Mesos cluster state
type Framework struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Tasks []Task `json:"tasks"`
}

// Task represents task info from Mesos cluster state
type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	SlaveID   string    `json:"slave_id"`
	State     string    `json:"state"`
	Resources Resources `json:"resources"`
}

// Resources represents resources from Mesos
type Resources struct {
	CPUs   float64 `json:"cpus"`
	Memory float64 `json:"mem"`
}

func (r Resources) String() string {
	return fmt.Sprintf("%.2f CPUs, %.2fGB RAM", r.CPUs, r.Memory/1024.)
}

// Add combines resources and updates itself with the new value
func (r *Resources) Add(more Resources) {
	r.CPUs += more.CPUs
	r.Memory += more.Memory
}

// Slave represents slave info from Mesos cluster state
type Slave struct {
	ID         string                 `json:"id"`
	Hostname   string                 `json:"hostname"`
	Attributes map[string]interface{} `json:"attributes"`
	// Total available resources
	Resources Resources `json:"resources"`
	// Reserved resources per role
	ReservedResources map[string]Resources `json:"reserved_resources"`
	// Unreserved resources (role *)
	UnreservedResources Resources `json:"unreserved_resources"`
}

// GetState returns Mesos master state from requested URL
func GetState(u *url.URL) (*State, error) {
	u.Path = "/master/state"
	r, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body from mesos: %s", err)
		}
	}()

	s := &State{}
	return s, json.NewDecoder(r.Body).Decode(s)
}
