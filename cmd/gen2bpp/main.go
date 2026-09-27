package main

import (
	"fmt"
	_ "image/png"
	"log"
)

func main() {
	ti, err := NewTilesImage("assets/dog.png")
	if err != nil {
		log.Fatal(err)
	}
	for i, tile := range ti.Tiles() {
		fmt.Printf("i: %d, data: %04x\n", i, EncodeTile(tile))
	}
}
