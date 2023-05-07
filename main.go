package main

import (
	"log"
	"os"
	"synchro-db/cmd"
)

func main() {

	cli := cmd.Setup()
	err := cli.Run(os.Args)
	if err != nil {
		log.Panicln("[-] error launching the app")
	}

	return
}
