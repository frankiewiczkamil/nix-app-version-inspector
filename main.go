package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var NIX_STORE_DEFAULT_PATH = "/nix/store"

func findActiveVersions(queriedName string) bytes.Buffer {
	// consider filtering out derivations (| grep -v '\.drv$')

	cmd1 := exec.Command("ls", NIX_STORE_DEFAULT_PATH, "-1")
	cmd2 := exec.Command("grep", queriedName)
	cmd3 := exec.Command("sed", "s|^|"+NIX_STORE_DEFAULT_PATH+"/|")

	pipe1, _ := cmd1.StdoutPipe()
	cmd2.Stdin = pipe1

	pipe2, _ := cmd2.StdoutPipe()
	cmd3.Stdin = pipe2

	var output bytes.Buffer
	cmd3.Stdout = &output

	cmd1.Start()
	cmd2.Start()
	cmd3.Start()
	cmd1.Wait()
	cmd2.Wait()
	cmd3.Wait()

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
