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

func runDir(dir string, dst string) error {
	dirName := filepath.Base(dir)
	files, err := filepath.Glob(filepath.Join(dir, "*.png"))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		return nil
	}

	romPath := filepath.Join(dst, dirName+".rom")
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

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("エラー: 2つの引数: src,dstが必要です\n")
		fmt.Printf("例: `go run ./cmd/gen2bpp ./assets/src ./assets/dst`\n")
		os.Exit(1)
	}

	src := os.Args[1]
	dst := os.Args[2]

	if !isDir(dst) {
		fmt.Printf("エラー: dst '%s'は有効なディレクトリではありません", dst)
		os.Exit(1)
	}

	matches, err := filepath.Glob(filepath.Join(src, "*"))
	if err != nil {
		log.Fatal(err)
	}
	info, err := os.Stat(dst)
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsDir() {
		fmt.Printf("dst: %sはディレクトリではありません\n", dst)
		os.Exit(1)
	}

	for _, path := range matches {
		if !isDir(path) {
			continue
		}
		if err := runDir(path, dst); err != nil {
			log.Fatalf("error in dir %s: %v", path, err)
		}
	}
}
