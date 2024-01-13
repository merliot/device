package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed build.tmpl
var buildTemplate string

type BuildConfig struct {
	Model       string
	ModelStruct string
}

func main() {
	model := flag.String("model", "foo", "Specify the model name")
	target := flag.String("target", "bar", "Specify the target name")
	flag.Parse()

	buildConfig := BuildConfig{
		Model:       *model,
		ModelStruct: strings.Title(*model),
	}

	// Create the build.go file
	file, err := os.CreateTemp("", "build-*.go")
	if err != nil {
		fmt.Printf("Error creating build.go: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(file.Name())

	tmpl, err := template.New("build").Parse(buildTemplate)
	if err != nil {
		fmt.Printf("Error reading build.tmpl: %v\n", err)
		os.Exit(1)
	}

	err = tmpl.Execute(file, buildConfig)
	if err != nil {
		fmt.Printf("Error executing build.tmpl: %v\n", err)
		os.Exit(1)
	}

	// Build the uf2 file
	uf2Name := *model + "-" + *target + ".uf2"
	cmd := exec.Command("tinygo", "build", "-target", *target, "-o", uf2Name, "-stack-size", "8kb", file.Name())
	println(cmd.String())
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%w: %s", err, stdoutStderr)
		os.Exit(1)
	}
}
