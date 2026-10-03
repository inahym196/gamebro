package main

import (
	"fmt"
	"os"

	"github.com/inahym196/gamebro/tools/gen2bpp"
)

func main() {

	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("引数はsrc, dstの2つ必要です.")
	}

	srcDir, dstDir := args[0], args[1]
	if err := gen2bpp.Generate(srcDir, dstDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
