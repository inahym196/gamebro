package gen2bpp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	_ "image/png"
	"io"
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

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func Generate(srcDir, dst string) error {
	if !isDir(srcDir) {
		return fmt.Errorf("%sはディレクトリではありません\n", srcDir)
	}

	files, err := filepath.Glob(filepath.Join(srcDir, "*.png"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create rom file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, file := range files {
		if err := appendTileImage(file, w); err != nil {
			return err
		}
	}
	return nil
}
