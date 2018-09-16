package main

import (
	"fmt"
	"os"

	arukas "github.com/dopin/arukas-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please specify a uuid of an Arukas service")
		os.Exit(1)
	}

	uuid := os.Args[1]

	if uuid == "" {
		fmt.Fprintln(os.Stderr, "Please specify a uuid of an Arukas service")
	}

	client := arukas.NewClientWithOsExitOnErr()
	service, err := client.GetService(uuid)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("ID: ", service.ID)
	fmt.Println("Image: ", service.Image)
	fmt.Println("Instances: ", service.Instances)
	fmt.Println("Plan: ", service.ServicePlan.Category)
}
