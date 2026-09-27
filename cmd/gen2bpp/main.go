package main

import (
	"bufio"
	"encoding/binary"
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

	bw := bufio.NewWriter(f)
	defer bw.Flush()

	for i, tile := range ti.Tiles() {
		bpps := EncodeTile(tile)
		binary.Write(bw, binary.NativeEndian, bpps)
		fmt.Printf("TileID: %d, data: %04x\n", i, bpps)
	}
	return nil
}

func main() {
	if err := run("./assets/dog.png", "./assets/dog.rom"); err != nil {
		log.Fatal(err)
	}
}
