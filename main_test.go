package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/webbben/code-exec-microservice/docker"
	"github.com/webbben/code-exec-microservice/execute"
)

// python snippets
var pythonCodeBasicPrint = "print(\"Hello world!\")"
var pythonCodeBasicLoop = "def looptyLoop(n):\n  output = 0\n  for i in range(n):\n    output = output+1\n  return output\nprint(looptyLoop(1000))"

// bash snippets
var bashCodeBasicPrint = "echo \"Hello world!\""
var bashCodeBasicLoop = "OUTPUT=0\nfor i in {1..1000}\ndo\n  OUTPUT=$((OUTPUT+1))\ndone\necho $OUTPUT"

// go snippets
var golangCodeBasicPrint = `
package main

import "fmt"

func main() {
	fmt.Print("Hello world!")
}
`

var golangCodeBasicLoop = `
package main

import "fmt"

func main() {
	sum := 0
	for i := 0; i < 1000; i++ {
		sum++
	}
	fmt.Print(sum)
}
`

// == Benchmarks ==

func BenchmarkPythonBasicLoop(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"lang": "python", "code": %s}`, pythonCodeBasicLoop))
	runBenchmarkExecRequest(b, reqData)
}

func BenchmarkPythonBasicPrint(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"lang": "python", "code": %s}`, pythonCodeBasicPrint))
	runBenchmarkExecRequest(b, reqData)
}

func BenchmarkBashBasicPrint(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"lang": "bash", "code": %s}`, bashCodeBasicPrint))
	runBenchmarkExecRequest(b, reqData)
}

func BenchmarkBashBasicLoop(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"lang": "bash", "code": %s}`, bashCodeBasicLoop))
	runBenchmarkExecRequest(b, reqData)
}

func BenchmarkGolangBasicPrint(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"lang": "go", "code": %s}`, golangCodeBasicPrint))
	runBenchmarkExecRequest(b, reqData)
}

func BenchmarkGolangBasicLoop(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"lang": "go", "code": %s}`, golangCodeBasicLoop))
	runBenchmarkExecRequest(b, reqData)
}

func runBenchmarkExecRequest(b *testing.B, requestData []byte) {
	server := httptest.NewServer(http.HandlerFunc(handleExecRequest))
	defer server.Close()

	client := server.Client()

	for i := 0; i < b.N; i++ {
		// Create the request using the input requestData
		req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(requestData))
		if err != nil {
			b.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Do the HTTP request
		_, err = client.Do(req)
		if err != nil {
			b.Fatalf("Error making request: %v", err)
		}
	}
}

// == Unit tests ==

func TestExecuteCodePythonBasic(t *testing.T) {
	var testCases = []struct {
		lang     string
		code     string
		expected string
	}{
		{lang: "python", code: pythonCodeBasicPrint, expected: "Hello world!"},
		{lang: "python", code: pythonCodeBasicLoop, expected: "1000"},
	}

	docker.InitDockerClient()

	for i, test := range testCases {
		testName := fmt.Sprintf("Python Basic %v", i)
		t.Run(testName, func(t *testing.T) {
			output, err := execute.ExecuteCode(test.lang, test.code, fmt.Sprintf("test%v", i))
			if err != nil {
				t.Errorf("ExecuteCode returned an error: %s", err.Error())
				return
			}
			if output != test.expected {
				t.Errorf("Output: [%s] Expected: [%s]", output, test.expected)
			}
		})
	}
}

func TestExecuteCodeBashBasic(t *testing.T) {
	var testCases = []struct {
		lang     string
		code     string
		expected string
	}{
		{lang: "bash", code: bashCodeBasicPrint, expected: "Hello world!"},
		{lang: "bash", code: bashCodeBasicLoop, expected: "1000"},
	}

	docker.InitDockerClient()

	for i, test := range testCases {
		testName := fmt.Sprintf("Bash Basic %v", i)
		t.Run(testName, func(t *testing.T) {
			output, err := execute.ExecuteCode(test.lang, test.code, fmt.Sprintf("test%v", i))
			if err != nil {
				t.Errorf("ExecuteCode returned an error: %s", err.Error())
				return
			}
			if output != test.expected {
				t.Errorf("Output: [%s] Expected: [%s]", output, test.expected)
			}
		})
	}
}

func TestExecuteCodeGolangBasic(t *testing.T) {
	var testCases = []struct {
		lang     string
		code     string
		expected string
	}{
		{lang: "go", code: golangCodeBasicPrint, expected: "Hello world!"},
		{lang: "go", code: golangCodeBasicLoop, expected: "1000"},
	}

	docker.InitDockerClient()

	for i, test := range testCases {
		testName := fmt.Sprintf("Golang Basic %v", i)
		t.Run(testName, func(t *testing.T) {
			output, err := execute.ExecuteCode(test.lang, test.code, fmt.Sprintf("test%v", i))
			if err != nil {
				t.Errorf("ExecuteCode returned an error: %s", err.Error())
				return
			}
			if output != test.expected {
				t.Errorf("Output: [%s] Expected: [%s]", output, test.expected)
			}
		})
	}
}
