package pyvalidate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Result struct {
	Name  *string
	Email *string
	Age   *int
}

type cliResponse struct {
	OK     bool           `json:"ok"`
	Detail string         `json:"detail"`
	Data   map[string]any `json:"data"`
}

// Validate runs the Python validation CLI with the given fields.
// Only keys present in fields are validated and returned.
func Validate(fields map[string]any) (*Result, error) {
	payload, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	python := pythonBin()
	script := scriptPath()

	cmd := exec.Command(python, script)
	cmd.Stdin = bytes.NewReader(payload)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	var resp cliResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		if runErr != nil {
			return nil, fmt.Errorf("python validator failed: %v: %s", runErr, stderr.String())
		}
		return nil, fmt.Errorf("invalid validator response: %w", err)
	}

	if !resp.OK {
		detail := resp.Detail
		if detail == "" {
			detail = "Validation failed"
		}
		return nil, &ValidationError{Detail: detail}
	}

	result := &Result{}
	if value, ok := resp.Data["name"].(string); ok {
		result.Name = &value
	}
	if value, ok := resp.Data["email"].(string); ok {
		result.Email = &value
	}
	if value, ok := resp.Data["age"].(float64); ok {
		age := int(value)
		result.Age = &age
	}

	return result, nil
}

type ValidationError struct {
	Detail string
}

func (e *ValidationError) Error() string {
	return e.Detail
}

func pythonBin() string {
	if value := os.Getenv("PYTHON_BIN"); value != "" {
		return value
	}
	return "python"
}

func scriptPath() string {
	if value := os.Getenv("PYTHON_VALIDATOR"); value != "" {
		return value
	}
	return filepath.Join("python", "validate_cli.py")
}
