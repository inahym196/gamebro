package main

import (
	"bufio"
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
		var buf [2]byte
		for _, bpp2 := range bpps {
			buf[0] = byte(bpp2 >> 8)
			buf[1] = byte(bpp2)
			bw.Write(buf[:])
		}
		fmt.Printf("TileID: %d, data: %04x\n", i, bpps)
	}
	return nil
}

func main() {
	if err := run("./assets/dog.png", "./assets/dog.rom"); err != nil {
		log.Fatal(err)
	}
}
