package dockerClient

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func ClientElement() {
	apiClient, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), client.ContainerListOptions{})

	if err != nil {
		panic(err)
	}

	for _, ctr := range containers.Items {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
}
