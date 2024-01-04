package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
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

var golangTestLoop = `
package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("main called in go script")
	sum := 0
	for i := 0; i < 1000; i++ {
		sum++
	}
	log.Println("main ended in go script")
	fmt.Print(sum)
}
`

// == Benchmarks ==

func BenchmarkPythonBasicLoop(b *testing.B) {
	runBenchmarkExecuteCode(b, "python", pythonCodeBasicLoop)
}

func BenchmarkPythonBasicPrint(b *testing.B) {
	runBenchmarkExecuteCode(b, "python", pythonCodeBasicPrint)
}

func BenchmarkBashBasicPrint(b *testing.B) {
	runBenchmarkExecuteCode(b, "bash", bashCodeBasicPrint)
}

func BenchmarkBashBasicLoop(b *testing.B) {
	runBenchmarkExecuteCode(b, "bash", bashCodeBasicLoop)
}

func BenchmarkGolangBasicPrint(b *testing.B) {
	runBenchmarkExecuteCode(b, "go", golangCodeBasicPrint)
}

func BenchmarkGolangBasicLoop(b *testing.B) {
	runBenchmarkExecuteCode(b, "go", golangCodeBasicLoop)
}

func runBenchmarkExecuteCode(b *testing.B, lang string, code string) {
	docker.InitDockerClient()
	for i := 0; i < b.N; i++ {
		execute.ExecuteCode(lang, code, fmt.Sprintf("benchmark%v", i), false)
	}
}

func BenchmarkParallelPython(b *testing.B) {
	runParallelBenchmark(b, "python", pythonCodeBasicLoop)
}

func BenchmarkParallelBash(b *testing.B) {
	runParallelBenchmark(b, "bash", bashCodeBasicLoop)
}

func BenchmarkParallelGo(b *testing.B) {
	runParallelBenchmark(b, "go", golangCodeBasicLoop)
}

func runParallelBenchmark(b *testing.B, lang string, code string) {
	docker.InitDockerClient()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			jobID := fmt.Sprintf("benchmark_%s", uuid.New().String())
			_, err := execute.ExecuteCode(lang, code, jobID, false)
			if err != nil {
				fmt.Printf("Error executing code: %v\n", err)
			}
		}
	})
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
			output, err := execute.ExecuteCode(test.lang, test.code, fmt.Sprintf("test%v", i), false)
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
			output, err := execute.ExecuteCode(test.lang, test.code, fmt.Sprintf("test%v", i), false)
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
			output, err := execute.ExecuteCode(test.lang, test.code, fmt.Sprintf("test%v", i), false)
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
