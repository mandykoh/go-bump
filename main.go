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
	const DefaultVersion = "0.0.0"

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go-bump <package_name> [semver]\n")
		os.Exit(1)
	}
	packageName := os.Args[1]

	tagPattern := regexp.MustCompile(`^\s*([^.]+)\.([^.]+)\.(.+)\s*$`)

	var revString, verMajor, verMinor, verRev string

	if len(os.Args) > 2 {
		revString = os.Args[2]

	} else {
		cmd := exec.Command("git", "describe", "--tags")
		rev, err := cmd.Output()
		if err != nil {
			revString = DefaultVersion
		} else {
			revString = string(rev)
		}
	}

	matches := tagPattern.FindStringSubmatch(revString)
	if matches == nil {
		fmt.Fprintf(os.Stderr, "Unable to parse revision '%s'\n", revString)
		os.Exit(1)
	}

	verMajor, verMinor, verRev = matches[1], matches[2], bumpVersion(matches[3])

	verFile, err := os.Create("version.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer verFile.Close()

	writer := bufio.NewWriter(verFile)
	writer.WriteString(fmt.Sprintf("package %s\n", packageName))
	writer.WriteString("\n")
	writer.WriteString(fmt.Sprintf("//go:generate go run vendor/github.com/mandykoh/go-bump/main.go %s\n", packageName))
	writer.WriteString("\n")
	writer.WriteString(fmt.Sprintf("const VersionMajor = \"%s\"\n", verMajor))
	writer.WriteString(fmt.Sprintf("const VersionMinor = \"%s\"\n", verMinor))
	writer.WriteString(fmt.Sprintf("const VersionRevision = \"%s\"\n", verRev))
	writer.WriteString(fmt.Sprintf("const Version = \"%s.%s.%s\"\n", verMajor, verMinor, verRev))
	writer.Flush()

	fmt.Printf("Generated version %s.%s.%s for package %s\n", verMajor, verMinor, verRev, packageName)
}
