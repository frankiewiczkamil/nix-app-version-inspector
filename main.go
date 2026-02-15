package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func findActiveVersions(queriedName string) bytes.Buffer {
	// consider filtering out derivations (| grep -v '\.drv$')
	// consider using pure ls for performance reasons

	cmd1 := exec.Command("nix-store", "--gc", "--print-live")
	cmd2 := exec.Command("grep", queriedName)

	pipe, _ := cmd1.StdoutPipe()
	cmd2.Stdin = pipe

	var output bytes.Buffer
	cmd2.Stdout = &output

	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	cmd2.Wait()

	return output
}

func getRootByPath(path string) ([]byte, error) {
	cmd := exec.Command("nix-store", "--query", "--roots", path)
	return cmd.Output()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument")
	}
	matchingVersionPaths := findActiveVersions(os.Args[1])

	lines := strings.SplitSeq(matchingVersionPaths.String(), "\n")
	for l := range lines {
		fmt.Println(l)
		if l != "" {
			refs, err := getRootByPath(l)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("\treferenced by:")
			for ref := range strings.SplitSeq(string(refs), "\n") {
				if ref != "" {
					fmt.Println("\t - " + string(ref))
				}
			}
		}
	}
}
