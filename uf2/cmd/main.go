package main

import (
	"log"
	"os"

	"github.com/merliot/device/uf2"
)

func main() {
	file, err := os.Open("foo.uf2")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	blocks, err := uf2.ReadUF2File(file)
	if err != nil {
		log.Fatal(err)
	}

	/*
		// Modify UF2 file as needed
		uf2.ModifyUF2File(func(block *uf2.UF2Block) {
			// Modify block data or other properties as needed
		})
	*/

	// Write modified UF2 file
	outputFile, err := os.Create("modified.uf2")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	if err := uf2.WriteUF2File(outputFile, blocks); err != nil {
		log.Fatal(err)
	}
}
