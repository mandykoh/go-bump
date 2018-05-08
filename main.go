package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go-version <package_name>\n")
		os.Exit(1)
	}
	packageName := os.Args[1]

	tagPattern := regexp.MustCompile(`^\s*([^.]+)\.([^.]+)\.(.+)\s*$`)

	cmd := exec.Command("git", "describe", "--tags")
	rev, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	matches := tagPattern.FindStringSubmatch(string(rev))
	if matches == nil {
		fmt.Fprintf(os.Stderr, "Unable to parse revision '%s'\n", string(rev))
		os.Exit(1)
	}

	verFile, err := os.Create("version.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer verFile.Close()

	writer := bufio.NewWriter(verFile)
	writer.WriteString(fmt.Sprintf("package %s\n", packageName))
	writer.WriteString("\n")
	writer.WriteString("//go:generate go run vendor/github.com/mandykoh/go-version/main.go\n")
	writer.WriteString("\n")
	writer.WriteString(fmt.Sprintf("const VersionMajor = \"%s\"\n", matches[1]))
	writer.WriteString(fmt.Sprintf("const VersionMinor = \"%s\"\n", matches[2]))
	writer.WriteString(fmt.Sprintf("const VersionRevision = \"%s\"\n", matches[3]))
	writer.WriteString(fmt.Sprintf("const Version = \"%s.%s.%s\"\n", matches[1], matches[2], matches[3]))
	writer.Flush()

	fmt.Printf("Generated version %s.%s.%s for package %s\n", matches[1], matches[2], matches[3], packageName)
}
