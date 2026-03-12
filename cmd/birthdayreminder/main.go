package main

import (
	"birthdayreminder/cmd/cli"
	"os"
)

func main() {
	cli, err := cli.NewCLI(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if err := cli.Run(); err != nil {
		panic(err)
	}
	cli.KillBot()
}
