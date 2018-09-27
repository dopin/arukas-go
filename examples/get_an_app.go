package main

import (
	"fmt"
	"os"

	arukas "github.com/dopin/arukas-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please specify a uuid of an Arukas app")
		os.Exit(1)
	}

	uuid := os.Args[1]

	if uuid == "" {
		fmt.Fprintln(os.Stderr, "Please specify a uuid of an Arukas app")
	}

	client := arukas.NewClientWithOsExitOnErr()
	app, err := client.GetApp(uuid)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Name: ", app.Name)
	fmt.Println("Image: ", app.Services[0].Image)
}
