package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/bobrik/scrappy/mesos"
	"github.com/bobrik/scrappy/report"
)

func main() {
	u := flag.String("u", "", "mesos url (http://host:port)")
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

	roles := map[string]*report.Role{}
	names := []string{}

	for _, agent := range rep.Agents {
		for _, role := range agent.Roles {
			if _, ok := roles[role.Name]; !ok {
				roles[role.Name] = &report.Role{
					Name:               role.Name,
					AvailableResources: mesos.Resources{},
					AllocatedResources: mesos.Resources{},
				}

				names = append(names, role.Name)
			}

			roles[role.Name].AvailableResources.Add(role.AvailableResources)
			roles[role.Name].AllocatedResources.Add(role.AllocatedResources)
		}
	}

	sort.Strings(names)

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 10, 0, 1, ' ', tabwriter.AlignRight)

	w.Write([]byte("role\tCPUs used\tCPUs total\tCPU %\tRAM used\tRAM total\tRAM %\t\n"))

	for _, name := range names {
		role := roles[name]
		fmt.Fprintf(
			w,
			"%s\t%.2f\t%.2f\t%.2f%%\t%.2fGB\t%.2fGB\t%.2f%%\t\n",
			role.Name,
			role.AllocatedResources.CPUs,
			role.AvailableResources.CPUs,
			role.AllocatedResources.CPUs/role.AvailableResources.CPUs*100,
			role.AllocatedResources.Memory/1024,
			role.AvailableResources.Memory/1024,
			role.AllocatedResources.Memory/role.AvailableResources.Memory*100,
		)
	}

	w.Flush()
}
