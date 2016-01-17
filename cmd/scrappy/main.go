package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bobrik/scrappy/mesos"
	"github.com/bobrik/scrappy/report"
)

func main() {
	u := flag.String("u", "", "mesos url (http://host:port)")
	s := flag.String("s", "host", "sort order for slaves: host, cpu, mem")
	r := flag.Bool("r", false, "reverse order")
	f := flag.String("f", "", "role name to filter on")

	flag.Parse()

	if *u == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	state, err := mesos.GetState(*u)
	if err != nil {
		log.Fatal(err)
	}

	rep := report.Generate(state, *f)
	report.SortSlaves(rep.Slaves, *s, *r)

	for i, slave := range rep.Slaves {
		fmt.Printf("%s: %s / %s (%v)\n", slave.Hostname, slave.AllocatedResources.String(), slave.AvailableResources.String(), slave.Attributes)

		fmt.Printf("  roles:\n")
		for _, role := range slave.Roles {
			fmt.Printf("    - %s: %s / %s\n", role.Name, role.AllocatedResources.String(), role.AvailableResources.String())
			fmt.Printf("      tasks: %d\n", len(role.Tasks))
			for _, task := range role.Tasks {
				fmt.Printf("        - %s: %s\n", task.Name, task.Resources.String())
			}
		}

		if i < len(rep.Slaves)-1 {
			fmt.Println()
		}
	}
}
