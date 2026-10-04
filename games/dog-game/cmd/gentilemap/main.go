package main

import (
	"fmt"
	"log"
	"os"
)

func random(x, y int) bool {
	n := uint32(x + y*20)

	n = (n ^ 61) ^ (n >> 16)
	n *= 9
	n = n ^ (n >> 4)
	n *= 0x27d4eb1d
	n = n ^ (n >> 15)

	return (n & 63) == 0
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Printf("引数は1つ dstが必要です")
		os.Exit(1)
	}

	romPath := args[0]
	width, height := 20, 18
	data := make([]byte, width*height)
	for y := range height {
		base := y * width
		for x := range width {
			if (x+y)%2 != 0 {
				data[base+x] = 0x7F
			} else if random(x, y) {
				data[base+x] = 0x01
			}
		}
	}
	if err := os.WriteFile(romPath, data, 0644); err != nil {
		log.Fatal(err)
	}
}
