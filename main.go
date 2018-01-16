package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
)

var (
	port = kingpin.Flag("port", "server listener port").Required().Int()
)

func main() {
	kingpin.Parse()
	fmt.Printf("%d\n", *port)
}
