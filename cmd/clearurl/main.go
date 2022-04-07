package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/intakeapp/clearurl"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("a parameter is required!")
	}

	h, err := clearurl.Init()
	fatalIfErr(err)
	var result []string
	for _, u := range strings.Split(args[0], "\n") {
		r, err := h.Clear(u)
		fatalIfErr(err)
		result = append(result, r)
	}
	fmt.Println(strings.Join(result, "\n"))
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
