package main

import (
	"fmt"
	"log"
	"os"
)

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
				data[base+x] = 1
			}
		}
	}
	if err := os.WriteFile(romPath, data, 0644); err != nil {
		log.Fatal(err)
	}
}
