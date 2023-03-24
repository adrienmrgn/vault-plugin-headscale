package headscale

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// containerConfig : Configuration to run a container using testcontainer
// Used as parameter input for runTestContainer func
type containerConfig struct {
	Image            string   // docker image name
	Name             string   // container name
	ConfigPath       string   // path to config file to be bind mount
	ConfigTargetPath string   // destination of mounted config file
	Command          []string // container command to be executed at startup
	Port             string   // Exposed port
}

// testContainer is returned by runTestContainer func
type testContainer struct {
	URI       string
	Container testcontainers.Container
	Context   context.Context
}

// Terminate wraps Terminate function on testContainer.Container
func (tc *testContainer) Terminate() error {
	return tc.Container.Terminate(tc.Context)
}

// retrieveHeadscaleConfigFromContainer
func retrieveHeadscaleConfigFromContainer(container string) (apiKey string, err error) {
	cmd := exec.Command("docker", "exec", container, "headscale", "apikey", "create", "-o", "yaml")
	apiKeyRaw, err := cmd.Output() // récupère la sortie de la commande
	if err != nil {
		return "", err
	}
	apiKey = strings.ReplaceAll(string(apiKeyRaw), "\n", "")
	return string(apiKey), nil
}

// runTestContainer
func runTestContainer(cc containerConfig) (container testContainer, err error) {

	mountSource := testcontainers.GenericBindMountSource{
		HostPath: cc.ConfigPath,
	}
	var mountTarget testcontainers.ContainerMountTarget = testcontainers.ContainerMountTarget(cc.ConfigTargetPath)

	req := testcontainers.ContainerRequest{
		Image:        cc.Image,
		Name:         cc.Name,
		Cmd:          cc.Command,
		ExposedPorts: []string{cc.Port + "/tcp"},
		WaitingFor:   wait.ForHTTP("/api/v1").WithPort(nat.Port(cc.Port + "/tcp")).WithStartupTimeout(time.Minute),
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source: mountSource,
				Target: mountTarget,
			},
		},
	}
	container.Context = context.Background()
	container.Container, err = testcontainers.GenericContainer(
		container.Context,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return container, err
	}

	mappedPort, err := container.Container.MappedPort(container.Context, nat.Port(cc.Port))
	if err != nil {
		return container, err
	}

	hostIP, err := container.Container.Host(container.Context)
	if err != nil {
		return container, err
	}

	container.URI = fmt.Sprintf("http://%s:%s", hostIP, mappedPort.Port())

	return container, err
}

func runHeadscale() (client *Client, container testContainer, err error) {
	pwd, _ := os.Getwd()
	headscale := containerConfig{
		Image:            "headscale/headscale:0.19",
		Name:             "headscale",
		Port:             "8080",
		ConfigPath:       path.Join(pwd, "test_config.yaml"),
		ConfigTargetPath: "/etc/headscale/config.yaml",
		Command:          []string{"headscale", "serve"},
	}
	container, err = runTestContainer(headscale)
	client = NewClient()
	client.APIURL = container.URI
	client.APIKey, err = retrieveHeadscaleConfigFromContainer(headscale.Name)
	return client, container, err
}
