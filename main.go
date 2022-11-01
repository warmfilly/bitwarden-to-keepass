package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/warmfilly/bitwarden-to-keepass/bw2kp"
)

func main() {
	opts := bw2kp.Options{}

	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	if err := bw2kp.GenerateKeepassDatabase(opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
