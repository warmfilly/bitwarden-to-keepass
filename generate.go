package main

import (
	"fmt"
	"os"
)

func GenerateKeepassDatabase(opts options) {
	fmt.Println(opts)

	file, err := os.Create(opts.DatabasePath)

	if err != nil {
		panic(err)
	}

	_ = file
}
