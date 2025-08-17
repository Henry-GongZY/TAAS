package main

import (
	cli2 "github.com/Henry-GongZY/TAAS/internal/cli"
	"os"
)

func main() {
	err := cli2.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
