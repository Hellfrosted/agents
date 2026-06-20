package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainHelp_printsShortHelp_whenInvokedWithHelp(t *testing.T) {
	// Given
	command := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "-h")
	command.Env = append(os.Environ(), "SK_UP_HELPER_PROCESS=1")

	// When
	output, err := command.CombinedOutput()

	// Then
	if err != nil {
		t.Fatalf("helper process failed: %v\n%s", err, output)
	}
	got := string(output)
	if !strings.Contains(got, "sk-up -I <source...>") {
		t.Fatalf("help output missing install-source shorthand:\n%s", got)
	}
}

func TestMainHelp_printsLongHelp_whenEntrypointOverrideProvided(t *testing.T) {
	// Given
	command := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "--help")
	command.Env = append(os.Environ(), "SK_UP_HELPER_PROCESS=1", "SK_UP_ENTRYPOINT=skills-updates")

	// When
	output, err := command.CombinedOutput()

	// Then
	if err != nil {
		t.Fatalf("helper process failed: %v\n%s", err, output)
	}
	got := string(output)
	if !strings.Contains(got, "skills-updates --install-source <source...>") {
		t.Fatalf("help output missing long install-source command:\n%s", got)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("SK_UP_HELPER_PROCESS") != "1" {
		return
	}
	os.Args = append([]string{"sk-up"}, os.Args[3:]...)
	main()
}
