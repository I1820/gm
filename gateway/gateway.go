/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-11-2017
 * |
 * | File Name:     gateway/gateway.go
 * +===============================================
 */

package runner

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/docker/go-connections/nat"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
)

var dockerClient *client.Client

// Gateway represents gateway docker information
type Gateway struct {
	ID   string `json:"id"`
	Port string `json:"port"`
}

func init() {
	// NewEnvClient initializes a new API client based on environment variables.
	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	dockerClient = cli
}

// New creates gateway docker with given user name
func New(name string) (Gateway, error) {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	ctx := context.Background()

	imageName := "loraserver/loraserver"

	_, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return Gateway{}, err
	}

	lport, _ := nat.NewPort("tcp", "8000")
	eport := fmt.Sprintf("%d", 8072+rand.Intn(100))

	resp, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				lport: struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				lport: []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: eport,
					},
				},
			},
		}, nil, fmt.Sprintf("da-%s", name))
	if err != nil {
		return Gateway{}, err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return Gateway{}, err
	}

	return Gateway{
		ID:   resp.ID,
		Port: eport,
	}, nil
}

// Remove removes runner docker
func (g Gateway) Remove() error {
	ctx := context.Background()

	err := dockerClient.ContainerRemove(ctx, g.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}
