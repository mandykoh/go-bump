package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func bumpVersion(versionPart string) string {
	ver, err := strconv.ParseUint(versionPart, 10, 64)
	if err != nil {
		return "0"
	}

	return strconv.FormatUint(ver + 1, 10)
}

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: go-version <bump> <package_name>\n")
		os.Exit(1)
	}
	bump := os.Args[1]
	packageName := os.Args[2]

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

	verMajor, verMinor, verRev := matches[1], matches[2], matches[3]

	switch bump {
	case "major": verMajor, verMinor, verRev = bumpVersion(verMajor), "0", "0"
	case "minor": verMinor, verRev = bumpVersion(verMinor), "0"
	default: case "revision": verRev = bumpVersion(verRev)
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
	writer.WriteString(fmt.Sprintf("//go:generate go run vendor/github.com/mandykoh/go-version/main.go %s\n", packageName))
	writer.WriteString("\n")
	writer.WriteString(fmt.Sprintf("const VersionMajor = \"%s\"\n", verMajor))
	writer.WriteString(fmt.Sprintf("const VersionMinor = \"%s\"\n", verMinor))
	writer.WriteString(fmt.Sprintf("const VersionRevision = \"%s\"\n", verRev))
	writer.WriteString(fmt.Sprintf("const Version = \"%s.%s.%s\"\n", verMajor, verMinor, verRev))
	writer.Flush()

	fmt.Printf("Generated version %s.%s.%s for package %s\n", verMajor, verMinor, verRev, packageName)
}
