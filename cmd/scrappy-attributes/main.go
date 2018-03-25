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

	attributes := make(map[string]*report.Attribute)
	ids := []string{}

	for _, slave := range rep.Slaves {
		for k, v := range slave.Attributes {
			id := fmt.Sprintf("%s=%v", k, v)

			if _, ok := attributes[id]; !ok {
				attributes[id] = &report.Attribute{
					Key:                k,
					Value:              v,
					AvailableResources: mesos.Resources{},
					AllocatedResources: mesos.Resources{},
				}

				ids = append(ids, id)
			}

			attributes[id].AvailableResources.Add(slave.AvailableResources)
			attributes[id].AllocatedResources.Add(slave.AllocatedResources)
		}
	}

	sort.Strings(ids)

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 10, 0, 1, ' ', tabwriter.AlignRight)

	w.Write([]byte("attribute\tCPUs used\tCPUs total\tCPU %\tRAM used\tRAM total\tRAM %\t\n"))

	for _, id := range ids {
		attribute := attributes[id]
		fmt.Fprintf(
			w,
			"%s\t%.2f\t%.2f\t%.2f%%\t%.2fGB\t%.2fGB\t%.2f%%\t\n",
			id,
			attribute.AllocatedResources.CPUs,
			attribute.AvailableResources.CPUs,
			attribute.AllocatedResources.CPUs/attribute.AvailableResources.CPUs*100,
			attribute.AllocatedResources.Memory/1024,
			attribute.AvailableResources.Memory/1024,
			attribute.AllocatedResources.Memory/attribute.AvailableResources.Memory*100,
		)
	}

	w.Flush()
}
