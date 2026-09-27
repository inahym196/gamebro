package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
)

func run(src, dst string) error {
	ti, err := NewTilesImage(src)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
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
	return nil
}

func main() {
	if err := run("./assets/dog.png", "./assets/dog.rom"); err != nil {
		log.Fatal(err)
	}
}
