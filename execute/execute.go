package execute

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/webbben/code-exec-microservice/docker"
)

var dockerImage = "py-golang:latest"

func CreateFile(lang string, code string, jobID string) (string, error) {
	var ext = ""
	switch lang {
	case "go":
		ext = "go"
	case "python":
		ext = "py"
	case "bash":
		ext = "sh"
		code = "#!/usr/bin/env bash\n" + code // add shebang - this should be portable to all distros?
	default:
		return "", errors.New(fmt.Sprintf("Language %s unsupported", lang))
	}
	filename := fmt.Sprintf("scripts/%s/generated_script.%s", jobID, ext)
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filename, []byte(code), 0644); err != nil {
		return "", err
	}
	return filename, nil
}

func ExecuteCode(lang string, code string, jobID string) (string, error) {
	fmt.Printf("executing %s code\n", lang)
	filename, err := CreateFile(lang, code, jobID)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error making file: %s", err.Error()))
	}

	// execute code in a new container
	output, err := docker.RunCodeContainer(jobID, lang, filename)
	if err != nil {
		return "", err
	}
	fmt.Println("... done!")
	outputStr := strings.TrimSpace(string(output))
	return outputStr, nil
}
