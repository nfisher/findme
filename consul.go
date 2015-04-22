package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

func registerCheck(client *api.Client, port int) {
	check := &api.AgentCheckRegistration{
		ID:    fmt.Sprintf("health%v", port),
		Name:  "simples health status",
		Notes: "TTL based health check.",
	}
	check.TTL = "30s"
	err := client.Agent().CheckRegister(check)
	if err != nil {
		log.Fatal(err)
	}
}

// registers to local consul node.
func register(client *api.Client, ip string, port int) {
	check := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("simples%v", port),
		Name:    "simples",
		Tags:    []string{"v1"},
		Address: ip,
		Port:    port,
		Check: &api.AgentServiceCheck{
			TTL: "30s",
		},
	}

	err := client.Agent().ServiceRegister(check)
	if err != nil {
		log.Fatal(err)
	}
}

func check(client *api.Client, port int) {
	check := fmt.Sprintf("service:simples%v", port)

	for {
		err := client.Agent().PassTTL(check, "Still running.")
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(28 * time.Second)
	}
}

func cleanup(client *api.Client, ip string, port int) {
	check := fmt.Sprintf("health%v", port)
	client.Agent().CheckDeregister(check)
	log.Println("Deregistered health check.")

	serviceId := fmt.Sprintf("simples%v", port)
	client.Agent().ServiceDeregister(serviceId)
	log.Println("Deregistered service.")
}
