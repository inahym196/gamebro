package main

import (
	"fmt"
	"log"
	"os"

	"github.com/inahym196/gamebro/tools/gen2bpp"
)

func main() {

	srcDir, dstDir := "./assets/src", "./assets/"
	if err := gen2bpp.Run(srcDir, dstDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// generate tilemap.rom
	romPath := "./assets/tilemap.rom"
	width, height := 20, 18
	data := make([]byte, width*height)
	for y := range height {
		base := y * width
		for x := range width {
			if (x+y)%2 != 0 {
				data[base+x] = 1
			}
		}
	}
	if err := os.WriteFile(romPath, data, 0644); err != nil {
		log.Fatal(err)
	}
}
