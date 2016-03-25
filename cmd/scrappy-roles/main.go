package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"

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

	for _, slave := range rep.Slaves {
		for _, role := range slave.Roles {
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

	for _, name := range names {
		role := roles[name]
		fmt.Printf("- %s: %s / %s\n", role.Name, role.AllocatedResources.String(), role.AvailableResources.String())
	}
}
