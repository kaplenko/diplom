// Command runner is the sandbox entrypoint that runs inside an isolated
// Docker container. It reads a JSON payload from stdin, compiles the
// user-submitted Go code, executes it against each test case, and
// prints a JSON verdict to stdout.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type input struct {
	Code           string     `json:"code"`
	TestCases      []testCase `json:"test_cases"`
	TimeoutSeconds int        `json:"timeout_seconds"`
}

type testCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type testResult struct {
	Index    int    `json:"index"`
	Passed   bool   `json:"passed"`
	Got      string `json:"got"`
	Expected string `json:"expected"`
}

type output struct {
	Status  string       `json:"status"`
	Score   int          `json:"score"`
	Stdout  string       `json:"stdout"`
	Stderr  string       `json:"stderr"`
	Details []testResult `json:"details,omitempty"`
}

func main() {
	var in input
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		fatal("failed to parse input: " + err.Error())
	}

	if in.TimeoutSeconds <= 0 {
		in.TimeoutSeconds = 10
	}
	perTestTimeout := time.Duration(in.TimeoutSeconds) * time.Second

	if err := os.MkdirAll("/tmp/solution", 0o755); err != nil {
		fatal("mkdir: " + err.Error())
	}
	if err := os.WriteFile("/tmp/solution/main.go", []byte(in.Code), 0o644); err != nil {
		fatal("write code: " + err.Error())
	}

	// Compile
	build := exec.Command("go", "build", "-o", "/tmp/solution/solution", "/tmp/solution/main.go")
	var buildStderr bytes.Buffer
	build.Stderr = &buildStderr
	if err := build.Run(); err != nil {
		writeOutput(output{
			Status: "failed",
			Stderr: buildStderr.String(),
		})
		return
	}

	passed := 0
	var details []testResult
	var summary strings.Builder

	for i, tc := range in.TestCases {
		var stdout, stderr bytes.Buffer
		cmd := exec.Command("/tmp/solution/solution")
		cmd.Stdin = strings.NewReader(tc.Input)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		done := make(chan error, 1)
		go func() { done <- cmd.Run() }()

		tr := testResult{Index: i + 1, Expected: strings.TrimSpace(tc.Expected)}

		select {
		case err := <-done:
			got := strings.TrimSpace(stdout.String())
			tr.Got = got

			if err != nil {
				tr.Got = strings.TrimSpace(stderr.String())
			} else if got == tr.Expected {
				tr.Passed = true
				passed++
			}

		case <-time.After(perTestTimeout):
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			tr.Got = "timeout"
		}

		details = append(details, tr)
		fmt.Fprintf(&summary, "Test %d: %v\n", i+1, tr.Passed)
	}

	total := len(in.TestCases)
	score := 0
	if total > 0 {
		score = (passed * 100) / total
	}

	status := "failed"
	if passed == total && total > 0 {
		status = "passed"
	}

	writeOutput(output{
		Status:  status,
		Score:   score,
		Stdout:  summary.String(),
		Details: details,
	})
}

func writeOutput(o output) {
	_ = json.NewEncoder(os.Stdout).Encode(o)
}

func fatal(msg string) {
	writeOutput(output{Status: "failed", Stderr: msg})
	os.Exit(1)
}
