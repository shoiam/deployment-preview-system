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

	containers, containerErr := apiClient.ContainerList(context.Background(), client.ContainerListOptions{})
	images, imageErr := apiClient.ImageList(context.Background(), client.ImageListOptions{All: true})

	if containerErr != nil {
		panic(containerErr)
	} else {
		fmt.Println("Continer fetch successful")
	}

	if imageErr != nil {
		panic(imageErr)
	} else {
		fmt.Println("Image fetch successful")
	}
	i := 1
	for _, ctr := range containers.Items {
		fmt.Printf("Container %d\n", i)
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
		i += 1
	}
	j := 1
	for _, image := range images.Items {
		fmt.Printf("Image %d\n", j)
		fmt.Printf("%s %s (status: %s)\n", image.ID, image.ParentID, image.Labels[image.Descriptor.Digest.Algorithm().String()])
		j += 1
	}
}
