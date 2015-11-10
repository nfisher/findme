package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

// registers to local consul node.
func register(client *api.Client, ip string, port int) {
	check := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("simples%v", port),
		Name:    "simples",
		Tags:    []string{"v1"},
		Address: ip,
		Port:    port,
		Checks: api.AgentServiceChecks{
			&api.AgentServiceCheck{
				TTL: "30s",
			},
			&api.AgentServiceCheck{
				Interval: "15s",
				HTTP:     fmt.Sprintf("http://127.0.0.1:%v/health", port),
			},
		},
	}

	err := client.Agent().ServiceRegister(check)
	if err != nil {
		log.Fatal("ServiceRegister", err)
	}
}

func check(client *api.Client, port int) {
	check := fmt.Sprintf("service:simples%v:1", port)

	for {
		err := client.Agent().PassTTL(check, "Still running.")
		if err != nil {
			log.Fatal("Check", err)
		}

		time.Sleep(28 * time.Second)
	}
}

func cleanup(client *api.Client, ip string, port int) {
	check := fmt.Sprintf("health%v", port)
	client.Agent().CheckDeregister(check)
	log.Println("Deregistered health checks")

	serviceId := fmt.Sprintf("simples%v", port)
	client.Agent().ServiceDeregister(serviceId)
	log.Println("Deregistered service.")
}
