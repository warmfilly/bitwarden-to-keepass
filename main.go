package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	opts := options{}

	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	if err := GenerateKeepassDatabase(opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
