package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
)

func appendTileImage(file string, w io.Writer) error {
	ti, err := NewTilesImage(file)
	if err != nil {
		return err
	}

	for i, tile := range ti.Tiles() {
		bpps := EncodeTile(tile)
		binary.Write(w, binary.NativeEndian, bpps)
		fmt.Printf("file: %s, TileID: %d, data: %04x\n", file, i, bpps)
	}
	return nil
}

func runDir(dir string) error {
	dirName := filepath.Base(dir)
	files, err := filepath.Glob(filepath.Join(dir, "*.png"))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		return nil
	}

	romPath := filepath.Join("assets/", dirName+".rom")
	f, err := os.Create(romPath)
	if err != nil {
		return fmt.Errorf("failed to create rom file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, file := range files {
		if err := appendTileImage(file, w); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func main() {
	dirs, err := filepath.Glob("./assets/src/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		if err := runDir(dir); err != nil {
			log.Fatalf("error in dir %s: %v", dir, err)
		}
	}
}
