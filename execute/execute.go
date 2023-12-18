package execute

import (
	"errors"
	"fmt"
	"os/exec"
)

var dockerImage = "py-golang:latest"

func ExecuteCode(lang string, code string) (string, error) {
	var cmd *exec.Cmd

	// execute code in a new container
	switch lang {
	case "go":
		fmt.Println("using golang")
		cmd = exec.Command("docker", "run", "--rm", "--net=none", dockerImage, "go", "run", code)
	case "python":
		fmt.Println("using python")
		cmd = exec.Command("docker", "run", "--rm", "--net=none", dockerImage, "python3", "-c", code)
	default: // defaults to bash
		fmt.Println("using bash")
		cmd = exec.Command("docker", "run", "--net=none", dockerImage, code)
	}

	// Capture the output of the container
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	if output == nil {
		return "", errors.New("Output was nil")
	}
	fmt.Printf("Raw Output: %s\n", output)
	fmt.Println("... done!")
	return string(output), nil
}
