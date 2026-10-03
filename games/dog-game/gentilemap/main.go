package main

import (
	"log"
	"os"
)

func main() {
	romPath := "./games/dog-game/assets/tilemap.rom"

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
