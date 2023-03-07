package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

// func TestIntegration(t *testing.T) {
// 	cmd := exec.Command("/bin/bash", "-c", "sudo ../run_integration_test.sh")
// 	err := cmd.Run()
// 	if err != nil {
// 		t.Errorf("Failed to run integration test: %v", err)
// 	}
// }

func TestIntegration(t *testing.T) {
	cmd := exec.Command("/bin/bash", "../scripts/run_integration_test.sh")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)
	err := cmd.Run()
	if err != nil {
		t.Errorf("Failed to run integration test: %v\n%s", err, stderr.String())
	}
	fmt.Printf("%s", stdout.String())
}
