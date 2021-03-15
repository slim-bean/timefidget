//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Runs go mod download and then installs the binary.
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return nil
}

// Build a docker image for amd64
func BuildDockerAMD() error {
	if err := sh.RunV("docker", "build", "-t", "slimbean/fidgserver:latest", "-f", "cmd/fidgserver/Dockerfile", "."); err != nil {
		return err
	}
	return nil
}

// Build a docker image for arm64
func ARM64Image() error {
	if err := sh.RunV("docker", "build", "--build-arg", "TARGET_PLATFORM=linux/arm64", "--build-arg", "COMPILE_GOARCH=arm64", "-t", "slimbean/fidgserver-arm:latest", "-f", "cmd/fidgserver/Dockerfile", "."); err != nil {
		return err
	}
	return nil
}

func ARM64Push() error {
	mg.Deps(ARM64Image)
	if err := sh.RunV("docker", "push", "slimbean/fidgserver-arm:latest"); err != nil {
		return err
	}
	return nil
}
