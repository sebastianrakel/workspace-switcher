package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func executeString(command string) (string, error) {
	args := strings.Fields(command)
	return execute(args[0], args[1:], nil)
}

func execute(command string, args []string, inputLines []string) (string, error) {
	fmt.Printf("%s %s\n", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr

	if len(inputLines) == 0 {
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}

		return string(output), cmd.Err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()

	for _, entry := range inputLines {
		io.WriteString(stdin, fmt.Sprintf("%s\n", entry))
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func openRofi(entries []string) (string, error) {
	args := []string{
		"-p",
		"Activate Workspace: ",
		"-format",
		"i",
		"-dmenu",
	}

	result, err := execute("rofi", args, entries)
	if err != nil {
		return "", err
	}

	resultIndex, err := strconv.Atoi(result)
	if err != nil {
		return "", err
	}

	return entries[resultIndex], nil
}
