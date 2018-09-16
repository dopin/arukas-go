package main

import (
	"fmt"

	arukas "github.com/dopin/arukas-go"
)

func main() {
	client := arukas.NewClientWithOsExitOnErr()
	if services, err := client.GetServices(); err == nil {
		for _, service := range services {
			fmt.Println(service.ID, service.Image, service.ServicePlan.ID)
		}
	} else {
		fmt.Println(err)
	}
}
