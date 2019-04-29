package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bobrik/scrappy/mesos"
	"github.com/bobrik/scrappy/report"
)

func main() {
	u := flag.String("u", "", "mesos url (http://host:port)")
	f := flag.String("f", "", "role name to filter on")
	c := flag.Float64("c", 0.01, "minimum CPU block to consider")
	m := flag.Float64("m", 32, "minimum memory block to consider, specified in MB")

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
	resource_offers := map[string]int64{}

	for _, slave := range rep.Slaves {
		for _, role := range slave.Roles {
			if _, ok := roles[role.Name]; !ok {
				roles[role.Name] = &report.Role{
					Name:               role.Name,
					AvailableResources: mesos.Resources{},
					AllocatedResources: mesos.Resources{},
				}
				resource_offers[role.Name] = 0

				names = append(names, role.Name)
			}

			roles[role.Name].AvailableResources.Add(role.AvailableResources)
			roles[role.Name].AllocatedResources.Add(role.AllocatedResources)
			cpu_offers := math.Floor((role.AvailableResources.CPUs - role.AllocatedResources.CPUs) / *c)
			mem_offers := math.Floor((role.AvailableResources.Memory - role.AllocatedResources.Memory) / *m)
			resource_offers[role.Name] += (int64)(math.Min(cpu_offers, mem_offers))
		}
	}

	sort.Strings(names)

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 10, 0, 1, ' ', tabwriter.AlignRight)

	// there should be a saner way to do this with fmt.sprintf....
	cpu_s := fmt.Sprintf("%f", *c)
	cpu_s = strings.TrimRight(strings.TrimRight(cpu_s, "0"), ".")
	w.Write([]byte(fmt.Sprintf("role\tCPUs used\tCPUs total\tCPU %%\tRAM used\tRAM total\tRAM %%\tOffers remaining for %.fMB, %s CPU\t\n", *m, cpu_s)))

	for _, name := range names {
		role := roles[name]
		fmt.Fprintf(
			w,
			"%s\t%.2f\t%.2f\t%.2f%%\t%.2fGB\t%.2fGB\t%.2f%%\t%d\t\n",
			role.Name,
			role.AllocatedResources.CPUs,
			role.AvailableResources.CPUs,
			role.AllocatedResources.CPUs/role.AvailableResources.CPUs*100,
			role.AllocatedResources.Memory/1024,
			role.AvailableResources.Memory/1024,
			role.AllocatedResources.Memory/role.AvailableResources.Memory*100,
			resource_offers[role.Name],
		)
	}

	w.Flush()
}
