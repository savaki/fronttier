package main

import (
	"flag"
	"github.com/savaki/fronttier"
)

var (
	port     = flag.Int("port", 8080, "the port to listen to [default:8080]")
	filename = flag.String("config", "fronttier.json", "the configuration file to read [default:fronttier.json]")
)

func getServer(filename string) (*fronttier.Router, error) {
	return fronttier.NewRouter(), nil
}

func main() {
	flag.Parse()

	// server, err := getServer(*filename)
	// if err != nil {
	// 	fmt.Printf("unable to start fronttier => %s\n", err)
	// 	os.Exit(1)
	// }

	// addr := fmt.Sprintf(":%d", *port)
	// panic(server.ListenAndServe(addr))
}
