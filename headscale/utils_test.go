package headscale

import (
	"os"
	"testing"
	"path"

	"github.com/stretchr/testify/assert"

)

func TestRunContainer(t *testing.T) {
	pwd, _ := os.Getwd()
	headscale := containerConfig{
		Image: "headscale/headscale:0.19",
		Name: "headscale",
		Port: "8080",
		ConfigPath: path.Join(pwd ,"test_config.yaml"),
		ConfigTargetPath: "/etc/headscale/config.yaml",
		Command: []string{"headscale","serve"},
	}
	container, err := runTestContainer(headscale)
	defer container.Container.Terminate(container.Context)
	assert.NoError(t, err)
	assert.Contains(t, container.URI, "http")
}