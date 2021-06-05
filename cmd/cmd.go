package cmd

import (
	"flag"
)

func Command() (string, string) {
	port := flag.String("port", "foo", "an int")
	host := flag.String("host", "localhost", "a string")

	flag.Parse()

	return *port, *host
}
