package main

import (
	"fmt"
	"os"

	arukas "github.com/dopin/arukas-go"
)

func main() {
	app := &arukas.App{
		Name: "My New App",
		Services: arukas.Services{
			&arukas.Service{
				ClientID:  "1",
				Image:     "nginx",
				Instances: 1,
				Ports: arukas.Ports{
					&arukas.Port{
						Protocol: "tcp",
						Number:   80,
					},
				},
				Environment: arukas.Environment{
					&arukas.Env{
						Key:   "FOO",
						Value: "bar",
					},
				},
				ServicePlan: &arukas.ServicePlan{
					ID: "jp-tokyo/hobby",
				},
			},
		},
	}

	client := arukas.NewClientWithOsExitOnErr()
	createdApp, err := client.CreateApp(app)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Created App ID: ", createdApp.ID)
}
