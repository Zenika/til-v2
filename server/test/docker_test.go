package main_test

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/ory/dockertest/v3/docker/opts"
	"testing"
	"time"
)

func spawnStack(t *testing.T) {
	var err error

	// Suffix will enable us to run concurrent tests instances and avoid container collision even if the destroyStack fails.
	suffix := fmt.Sprintf("%d", time.Now().UnixMilli())

	client, err := docker.NewVersionedClient(opts.DefaultHost, "1.52")
	if err != nil {
		destroyStack()
	}

	pool = &dockertest.Pool{
		Client: client,
	}

	// Create mock server
	mock, err = pool.BuildAndRunWithOptions("google_mock/Dockerfile", &dockertest.RunOptions{
		Name:         fmt.Sprintf("til-mock-authentication-server_%s", suffix),
		ExposedPorts: []string{"8000/tcp"},
	})
	if err != nil {
		destroyStack()
	}

	// Create main server with matching build args to use our mock server
	server, err = pool.BuildAndRunWithOptions("../Dockerfile", &dockertest.RunOptions{
		Name:         fmt.Sprintf("til-server_%s", suffix),
		ExposedPorts: []string{"8000/tcp"},
		Env: []string{
			"TIL_JWT_SECRET=azertyuiopqsdfghjklmwxcvbn",
			"TIL_GOOGLE_CLIENT_ID=placeholder",
			"TIL_GOOGLE_CLIENT_SECRET=placeholder",
			"TIL_DEBUG=1",
			"TIL_USE_IN_MEMORY_DATABASE=1",
			"TIL_DATABASE_FILE_NAME=:memory:",
			"TIL_DEFAULT_ADMIN=102950075881792615000",
			fmt.Sprintf("TIL_GOOGLE_TOKEN_ENDPOINT=http://172.17.0.1:%s/", mock.GetPort("8000/tcp")),
		},
	})

	if err != nil {
		destroyStack()
	}
}

// destroyStack will keep your docker clean by removing all containers and images related to our tests at the end of the script.
func destroyStack() {
	if pool == nil {
		return
	}

	containers := []*dockertest.Resource{mock, server}
	for _, container := range containers {
		if container != nil {
			_ = pool.Purge(container)
			_ = pool.Client.RemoveImage(container.Container.Image)
		}
	}

	_, _ = pool.Client.PruneImages(docker.PruneImagesOptions{})
}
