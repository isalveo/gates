package main

import (
	"flag"
	"fmt"
	"initializers"
	"net/http"

	_ "server"
)

var p = flag.Int("p", 3000, "The port to listen on")

func main() {
	flag.Parse()

	if err := initializers.Boot(initializers.ConfigPaths, "hermes-gatekeeper"); err != nil {
		panic(err)
	}

	fmt.Printf("Listen and serve on, %d\n", *p)
	if err := http.ListenAndServe(fmt.Sprint(":", *p), nil); err != nil {
		panic(err)
	}
}
