package main

import (
	_ "embed"
	"log"
)

func main() {
	if err := NewRootCommnad().Execute(); err != nil {
		log.Fatal(err)
	}
}
