package execute

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

func ExecuteCode(lang string, code string, debug bool) (string, error) {
	var cmd *exec.Cmd
	if debug {
		log.Printf("start: executing %s code\n", lang)
	}

	// put the code in a new file
	jobID := uuid.New().String()
	filename, err := CreateFile(lang, code, jobID)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error making file: %s", err.Error()))
	}
	defer func() {
		// delete created folder/file
		if err = os.RemoveAll(fmt.Sprintf("scripts/%s", jobID)); err != nil {
			log.Printf("Failed to remove scripts directory for job %s: %s\n", jobID, err)
		}
	}()

	// get the correct command to run
	switch lang {
	case "go":
		cmd = exec.Command("go", "run", filename)
	case "python":
		cmd = exec.Command("python3", filename)
	case "bash":
		cmd = exec.Command("/bin/bash", filename)
	}

	// Capture the output
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
	outputStr := strings.TrimSpace(string(output))
	if debug {
		log.Printf("end: executing %s code\n", lang)
	}
	return outputStr, nil
}
