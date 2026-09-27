package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
)

func main() {
	ti, err := NewTilesImage("assets/dog.png")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("./assets/dog.rom")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bw := NewBPP2Writer(f)
	defer bw.Flush()

	for i, tile := range ti.Tiles() {
		bpps := EncodeTile(tile)
		for _, bpp2 := range bpps {
			bw.Write(bpp2)
		}
		fmt.Printf("i: %d, data: %04x\n", i, bpps)
	}
}
