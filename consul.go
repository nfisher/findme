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
	check := &api.CatalogRegistration{
		Datacenter: "dc1",
		Node:       "le-bigmac.local",
		Address:    ip,
		Service: &api.AgentService{
			ID:      fmt.Sprintf("simples%v", port),
			Service: "simples",
			Tags:    []string{"v1"},
			Address: ip,
			Port:    port,
		},
		Check: &api.AgentCheck{
			CheckID:   fmt.Sprintf("health%v", port),
			ServiceID: fmt.Sprintf("simples%v", port),
		},
	}

	meta, err := client.Catalog().Register(check, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(meta.RequestTime)
}

func check(client *api.Client, port int) {
	check := fmt.Sprintf("health%v", port)

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

	cat := &api.CatalogDeregistration{
		Datacenter: "dc1",
		Node:       "le-bigmac.local",
		Address:    ip,
		ServiceID:  fmt.Sprintf("simples%v", port),
		CheckID:    fmt.Sprintf("health%v", port),
	}
	client.Catalog().Deregister(cat, nil)
}
