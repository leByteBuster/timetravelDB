package main

import (
	"os/exec"
	"strings"
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
	sb := new(strings.Builder)
	sbErr := new(strings.Builder)
	cmd.Stdout = sb
	cmd.Stderr = sbErr
	//cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
	//cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to run integration test: %v\n%s", err, sbErr.String())
	}

	// Wait for the command to complete and then print the output.
	// if err := cmd.Wait(); err != nil {
	// 	t.Errorf("Failed to wait for command: %v", err)
	// }

	t.Log(sb.String())

	//fmt.Printf("%s", stdout.String())
	//output, _ := cmd.CombinedOutput()
}
