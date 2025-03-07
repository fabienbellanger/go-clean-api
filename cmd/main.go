package main

import (
	"go-clean-api/pkg/infrastructure/cli"
	"log"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalln(err)
	}
}
