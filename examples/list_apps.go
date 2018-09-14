package main

import (
	"fmt"

	arukas "github.com/dopin/arukas-resources-go"
)

func main() {
	client := arukas.NewClientWithOsExitOnErr()
	if apps, err := client.GetApps(); err == nil {
		for _, app := range apps {
			fmt.Println(app.ID, app.Name, app.Services[0].Image, app.Services[0].ServicePlan.ID)
		}
	} else {
		fmt.Println(err)
	}
}
