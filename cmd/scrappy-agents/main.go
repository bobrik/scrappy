package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/bobrik/scrappy/mesos"
	"github.com/bobrik/scrappy/report"
)

func main() {
	u := flag.String("u", "", "mesos url (http://host:port)")
	s := flag.String("s", "host", "sort order for agents: host, cpu, cpu_percent, mem, mem_percent, tasks")
	r := flag.Bool("r", false, "reverse order")
	f := flag.String("f", "", "role name to filter on")

	flag.Parse()

	if *u == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	mesosUrl, err := url.Parse(*u)
	if err != nil {
		log.Fatal(err)
	}

	state, err := mesos.GetState(mesosUrl)
	if err != nil {
		log.Fatal(err)
	}

	rep := report.Generate(state, *f)
	report.SortAgents(rep.Agents, *s, *r)

	for i, agent := range rep.Agents {
		fmt.Printf("%s: %s / %s (%v)\n", agent.Hostname, agent.AllocatedResources.String(), agent.AvailableResources.String(), agent.Attributes)

		fmt.Printf("  roles:\n")
		for _, role := range agent.Roles {
			fmt.Printf("    - %s: %s / %s\n", role.Name, role.AllocatedResources.String(), role.AvailableResources.String())
			fmt.Printf("      tasks: %d\n", len(role.Tasks))
			for _, task := range role.Tasks {
				fmt.Printf("        - %s: %s\n", task.Name, task.Resources.String())
			}
		}

		if i < len(rep.Agents)-1 {
			fmt.Println()
		}
	}
}
