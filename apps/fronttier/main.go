package main

import (
	"flag"
	"fmt"
	. "github.com/savaki/fronttier"
	"github.com/savaki/fronttier/core"
	"os"
)

var (
	port     = flag.Int("port", 8080, "the port to listen to [default:8080]")
	filename = flag.String("config", "fronttier.json", "the configuration file to read [default:fronttier.json]")
)

func getServer(filename string) (*core.Frontier, error) {
	flag.Parse()

	builder := Builder().LoadFile(filename)
	return builder.Build()
}

func main() {
	server, err := getServer(*filename)
	if err != nil {
		fmt.Printf("unable to start fronttier => %s\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", *port)
	panic(server.ListenAndServe(addr))
}
