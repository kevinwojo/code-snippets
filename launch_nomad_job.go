package main

import (
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"os"
)

var (
	nomadToken = os.Getenv("NOMAD_TOKEN")
	nomadAddress = os.Getenv("NOMAD_ADDR")
)

func main() {

	client, err := nomad.NewClient(&nomad.Config{
		Address: nomadAddress,
		SecretID: nomadToken,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Some variables you can play around with
	// Depending on your situation, these may live
	// somewhere else in your code.
	regions, err := client.Regions().List()
	region := regions[0] // dc1 in my configuration
	datacenter := "dc1"
	namespace := "default"
	cpuLimit := 100
	memLimit := 128
	taskCount := 1
	taskName := "go-nomad-fun"
	fabioPrefix := "urlprefix-"

	// Fetch the jobs endpoint
	jobs := client.Jobs()

	// Define task
	task := &nomad.Task{
		Name:   "nomad-fun",
		Driver: "docker",
		Config: map[string]interface{}{
			"image": "nginx:latest",
			"port_map": []map[string]interface{}{
				map[string]interface{}{"http": 80},
			},
		},
		Services: []*nomad.Service{&nomad.Service{
			Name:        taskName,
			Tags:        []string{fabioPrefix+"/"+taskName + " strip=/"+taskName},
			PortLabel:   "http",
			AddressMode: "",
			Checks: []nomad.ServiceCheck{nomad.ServiceCheck{
				Name:      "tcp",
				Type:      "tcp",
				Protocol:  "tcp",
				PortLabel: "http",
				Interval:  30000000000,
				Timeout:   60000000000,
			},
			}},
		},
		Resources: &nomad.Resources{
			CPU:      &cpuLimit,
			MemoryMB: &memLimit,
			Networks: []*nomad.NetworkResource{&nomad.NetworkResource{
				DynamicPorts: []nomad.Port{nomad.Port{Label: "http"}},
			}},
		},
	}

	// Define the task group
	taskGroup := &nomad.TaskGroup{
		Name:             &taskName,
		Count:            &taskCount,
		Tasks:            []*nomad.Task{task},
	}
	//  Get a Job type with default values for a service job
	job := nomad.NewServiceJob(taskName, taskName, region, 100)
	job.TaskGroups = []*nomad.TaskGroup{taskGroup}
	job.Datacenters = []string{datacenter}

	// This runs the job
	resp,_,err := jobs.Register(job, &nomad.WriteOptions{
		Namespace: namespace,
		AuthToken: nomadToken,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Some information and confirmation
	fmt.Println("Job Launched. Evaluation ID: ", resp.EvalID)
	allocations, _, err := jobs.Allocations(taskName, true, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=============== Allocations =============== ")
	for i := 0; i < len(allocations); i++ {
		fmt.Println(allocations[i].NodeName + ": " + allocations[i].ID)
	}
}