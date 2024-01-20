package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// docker client and images

var (
	dockerClient *client.Client       // docker client
	image        = "py-golang:latest" // image our containers use
)

// resource limits for containers

var (
	maxMemory int64 = 128 * 1024 * 1024 // 128 MB
	maxCPU    int64 = int64(1e9)        // 1 CPU core
)

// other config

var (
	logStatus  bool = false // log container exit status?
	autoRemove bool = true  // automatically remove container after code execution?
)

// Initializes the docker client
func InitDockerClient() {
	var err error
	dockerClient, err = client.NewClientWithOpts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Docker client initialization complete")
}

func RunCodeContainer(jobID string, lang string, filename string, debug bool) (string, error) {
	ctx := context.Background()

	// create the command based on the programming language
	filePath := fmt.Sprintf("/app/%s", filename)
	var cmd = []string{}
	switch lang {
	case "go":
		cmd = []string{"go", "run", filePath}
	case "python":
		cmd = []string{"python3", filePath}
	case "bash":
		cmd = []string{"/bin/bash", filePath}
	default:
		return "", errors.New(fmt.Sprintf("language %s not supported", lang))
	}

	containerID, err := createContainer(jobID, cmd)
	if err != nil {
		return "", err
	}
	if debug {
		log.Println("container created")
	}
	// docker SDK's built-in autoRemove setting seems to cause problems for bash code execution, so doing this instead
	defer func() {
		if autoRemove {
			err = dockerClient.ContainerRemove(ctx, jobID, types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})
			if err != nil {
				fmt.Println("Error removing container:", err)
			}
		}
	}()

	// run the container
	if debug {
		log.Println("starting container")
	}
	err = dockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}
	statusCh, errCh := dockerClient.ContainerWait(ctx, containerID, "")
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case status := <-statusCh:
		if logStatus || debug {
			log.Printf("Container exited with status %v\n", status)
		}
	}

	// get the output from the container
	out, err := dockerClient.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: false})
	if err != nil {
		return "", err
	}
	if debug {
		log.Println("reading output bytes")
	}
	dstout, dsterr := &bytes.Buffer{}, &bytes.Buffer{}
	_, err = stdcopy.StdCopy(dstout, dsterr, out)
	if err != nil {
		return "", err
	}
	outputString := dstout.String()
	if outputString == "" {
		outputString = dsterr.String()
	}
	outputString = removeNonPrintableChars(outputString)
	return outputString, nil
}

// strips non-printable ascii characters from a string
//
// docker container logs prefix weird non-ascii characters to the output for some reason.
// it seems like it could be the timestamp but it doesn't come through as that when the
// output bytes are decoded.
func removeNonPrintableChars(logStr string) string {
	// regular expression pattern that matches non-printable ascii characters
	regexPattern := "[[:cntrl:]]"
	re := regexp.MustCompile(regexPattern)
	return re.ReplaceAllString(logStr, "")
}

// create a docker container for a code execution job
func createContainer(jobID string, cmd []string) (string, error) {
	ctx := context.Background()
	config := &container.Config{
		Image:        image,
		Cmd:          cmd,
		WorkingDir:   "/app",
		AttachStdout: true,
		AttachStderr: true,
	}
	// get absolute path our working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	wd = filepath.ToSlash(wd)
	hostConfig := &container.HostConfig{
		AutoRemove:  false,
		Binds:       []string{fmt.Sprintf("%s/scripts/%s:/app/scripts/%s", wd, jobID, jobID)},
		NetworkMode: "none",
		Resources: container.Resources{
			Memory:   maxMemory,
			NanoCPUs: maxCPU,
		},
	}
	// create the container
	resp, err := dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,
		nil,
		jobID,
	)
	return resp.ID, nil
}
