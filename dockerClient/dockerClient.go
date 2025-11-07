package dockerClient

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func ClientElement(branch string) (string, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer apiClient.Close()

	ctx := context.Background()

	containerName := "preview-" + branch

	resp, err := apiClient.ContainerCreate(ctx, &container.Config{
		Image: "nginx:alpine",
		ExposedPorts: nat.PortSet{
			"80/tcp": struct{}{},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{{HostPort: "0"}},
		},
	}, nil, nil, containerName)

	if err != nil {
		return "", err
	}
	if err := apiClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}
	inspect, _ := apiClient.ContainerInspect(ctx, resp.ID)
	port := inspect.NetworkSettings.Ports["80/tcp"][0].HostPort

	return fmt.Sprintf("http://localhost:%s", port), nil
}
