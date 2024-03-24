package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/merliot/device"
	"github.com/merliot/device/uf2"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: program uf2-file")
	}

	file := os.Args[1]

	uf2, err := uf2.Read(file)
	if err != nil {
		log.Fatal("Error reading/parsing uf2 file:", file, err.Error())
	}

	data := uf2.Bytes()
	magic := []byte(device.UF2Magic)
	mlen := len(magic)

	start := bytes.Index(data, magic)
	if start == -1 {
		log.Fatal("start UF2Magic sequence not found")
	}
	println("start", start)

	// Find the end of the chunk
	end := bytes.Index(data[start+mlen:], magic)
	if end == -1 {
		log.Fatal("end UF2Magic sequence not found")
	}
	println("end", end)

	fmt.Println(string(data[start-15 : start+mlen+end+mlen+2]))
}
