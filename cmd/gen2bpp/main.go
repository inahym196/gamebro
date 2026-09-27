package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
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
	dirs, err := filepath.Glob("./assets/src/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		dirName := filepath.Base(dir)
		romPath := filepath.Join("assets/", dirName+".rom")
		files, err := filepath.Glob(filepath.Join(dir, "*.png"))
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if err := run(file, romPath); err != nil {
				log.Fatal(err)
			}
		}
	}
}
