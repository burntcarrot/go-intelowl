package main

import (
	"log"

	"github.com/burntcarrot/gointelowl/command"
)

func main() {
	if err := command.Execute(); err != nil {
		log.Fatalln(err)
	}
}
