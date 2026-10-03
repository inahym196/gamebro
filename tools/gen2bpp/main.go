package gen2bpp

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

func Run(src, dst string) error {

	if !isDir(dst) {
		return fmt.Errorf("エラー: dst '%s'は有効なディレクトリではありません", dst)
	}

	matches, err := filepath.Glob(filepath.Join(src, "*"))
	if err != nil {
		return err
	}
	info, err := os.Stat(dst)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("dst: %sはディレクトリではありません\n", dst)
	}

	for _, path := range matches {
		if !isDir(path) {
			continue
		}
		if err := runDir(path, dst); err != nil {
			return fmt.Errorf("error in dir %s: %v", path, err)
		}
	}
	return nil
}
