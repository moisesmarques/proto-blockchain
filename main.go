package main

import (
	"os"

	"judge-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
