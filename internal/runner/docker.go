package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type DockerRunner struct {
	Image       string
	Timeout     time.Duration
	MemoryLimit string
	CPULimit    string
}

type RunInput struct {
	Code           string          `json:"code"`
	TestCases      json.RawMessage `json:"test_cases"`
	TimeoutSeconds int             `json:"timeout_seconds"`
}

type RunOutput struct {
	Status string `json:"status"`
	Score  int    `json:"score"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func NewDockerRunner(image string, timeoutSec int, memory, cpu string) *DockerRunner {
	return &DockerRunner{
		Image:       image,
		Timeout:     time.Duration(timeoutSec) * time.Second,
		MemoryLimit: memory,
		CPULimit:    cpu,
	}
}

func (r *DockerRunner) Run(ctx context.Context, code string, testCases json.RawMessage) (*RunOutput, error) {
	input := RunInput{
		Code:           code,
		TestCases:      testCases,
		TimeoutSeconds: int(r.Timeout.Seconds()),
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal runner input: %w", err)
	}

	// Add a buffer on top of the runner's internal timeout so
	// the container has time to start and the runner can clean up.
	ctx, cancel := context.WithTimeout(ctx, r.Timeout+30*time.Second)
	defer cancel()

	args := []string{
		"run", "--rm",
		"--network=none",
		"--memory=" + r.MemoryLimit,
		"--cpus=" + r.CPULimit,
		"--pids-limit=64",
		"-i",
		r.Image,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return &RunOutput{
			Status: "failed",
			Score:  0,
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		}, nil
	}

	var output RunOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return &RunOutput{
			Status: "failed",
			Score:  0,
			Stdout: stdout.String(),
			Stderr: fmt.Sprintf("failed to parse runner output: %v", err),
		}, nil
	}

	return &output, nil
}
