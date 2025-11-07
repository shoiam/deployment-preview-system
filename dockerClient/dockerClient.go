package dockerClient

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func cleanUpPreview(branch string) error {
	ctx := context.Background()
	cleanApi, _ := client.NewClientWithOpts(client.FromEnv)
	containerName := "preview-app-" + branch
	cleanApi.ContainerStop(ctx, containerName, container.StopOptions{})
	cleanApi.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
	return nil
}

func GetPreviews() ([]types.Container, error) {
	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	containers, err := apiClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func ClientElement(branch string) (string, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer apiClient.Close()

	ctx := context.Background()

	containerName := "preview-app-" + branch

	filterArgs := filters.NewArgs()
	filterArgs.Add("name", containerName)

	containers, err := apiClient.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		return "", err
	}

	// If container exists, clean it up before creating new one
	if len(containers) > 0 {
		for _, cont := range containers {
			apiClient.ContainerStop(ctx, cont.ID, container.StopOptions{})
			err = apiClient.ContainerRemove(ctx, cont.ID, container.RemoveOptions{Force: true})
			if err != nil {
				return "", fmt.Errorf("failed to remove existing container: %w", err)
			}
		}
	}

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
