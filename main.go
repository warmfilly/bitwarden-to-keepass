package main

import "github.com/jessevdk/go-flags"

func main() {
	opts := options{}

	args, err := flags.Parse(&opts)

	_ = args

	if err == nil {
		GenerateKeepassDatabase(opts)
	}
}
