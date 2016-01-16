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

	flag.Parse()

	if *u == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	state, err := mesos.GetState(*u)
	if err != nil {
		log.Fatal(err)
	}

	report := report.Generate(state)

	for i, slave := range report.Slaves {
		fmt.Printf("%s: %s / %s\n", slave.Hostname, slave.AllocatedResources.String(), slave.AvailableResources.String())

		fmt.Printf("  roles:\n")
		for _, role := range slave.Roles {
			fmt.Printf("    - %s: %s / %s\n", role.Name, role.AllocatedResources.String(), role.AvailableResources.String())
			fmt.Printf("      tasks: %d\n", len(role.Tasks))
			for _, task := range role.Tasks {
				fmt.Printf("        - %s: %s\n", task.Name, task.Resources.String())
			}
		}

		if i < len(report.Slaves)-1 {
			fmt.Println()
		}
	}
}
