package execute

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var dockerImage = "py-golang:latest"

func CreateFile(lang string, code string) (string, error) {
	var ext = ""
	switch lang {
	case "go":
		ext = "go"
	case "python":
		ext = "py"
	default:
		ext = "sh"
		code = "#!/usr/bin/env bash\n" + code // add shebang - this should be portable to all distros?
	}
	filename := fmt.Sprintf("scripts/generated_script.%s", ext)
	if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
		return "", err
	}
	return filename, nil
}

func ExecuteCode(lang string, code string) (string, error) {
	var cmd *exec.Cmd

	fmt.Printf("executing %s code\n", lang)
	filename, err := CreateFile(lang, code)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error making file: %s", err.Error()))
	}

	// execute code in a new container
	switch lang {
	case "go":
		// --rm: remove container on exit, --net=none: no internet access
		cmd = exec.Command("docker", "run", "--rm", "-v", "./scripts:/app/scripts", "--net=none", dockerImage, "go", "run", filename)
	case "python":
		cmd = exec.Command("docker", "run", "--rm", "-v", "./scripts:/app/scripts", "--net=none", dockerImage, "python3", filename)
	default: // defaults to bash
		cmd = exec.Command("docker", "run", "--rm", "-v", "./scripts:/app/scripts", "--net=none", dockerImage, "/bin/bash", filename)
	}

	// Capture the output of the container
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			err = errors.New(fmt.Sprintf("(%s) %s", err.Error(), output))
		}
		return "", err
	}
	if output == nil {
		return "", errors.New("Output was nil")
	}
	fmt.Println("... done!")
	outputStr := strings.TrimSpace(string(output))
	return outputStr, nil
}
