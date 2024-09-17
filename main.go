package main

import (
	"fmt"
	"log"
	"os"

	"gocli/cli"
)

func main() {
	opts, err := cli.ParseArgs(os.Args[1:], getEnvs())
	if err != nil {
		log.Fatalf("Error parsing arguments: %v", err)
	}

	fmt.Printf("Options: %+v\n", opts)
	// Proceed with using opts...
}

func getEnvs() map[string]string {
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		pair := splitOnce(e, '=')
		envMap[pair[0]] = pair[1]
	}
	return envMap
}

func splitOnce(s string, sep rune) [2]string {
	for i, c := range s {
		if c == sep {
			return [2]string{s[:i], s[i+1:]}
		}
	}
	return [2]string{s, ""}
}
