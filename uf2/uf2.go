package uf2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	UF2MagicStart0 = 0x0A324655 // "UF2\n"
	UF2MagicStart1 = 0x9E5D5157 // "Q\u009D"
	UF2BlockSize   = 512
)

// UF2File represents the structure of a UF2 file.
type UF2File struct {
	Blocks []UF2Block
}

// UF2Block represents a UF2 block.
type UF2Block struct {
	MagicStart0 uint32
	MagicStart1 uint32
	Flags       uint32
	TargetAddr  uint32
	PayloadSize uint32
	BlockNo     uint32
	NumBlocks   uint32
	FileSize    uint32 // or familyID;
	Data        [476]byte
	MagicEnd    uint32
}

// ReadUF2File reads a UF2 file from the given reader.
func ReadUF2File(reader io.Reader) (*UF2File, error) {
	var uf2 UF2File

	for {
		var block UF2Block
		err := binary.Read(reader, binary.LittleEndian, &block)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		uf2.Blocks = append(uf2.Blocks, block)
	}

	return &uf2, nil
}

// WriteUF2File writes a UF2 file to the given writer.
func WriteUF2File(writer io.Writer, uf2 *UF2File) error {
	for _, block := range uf2.Blocks {
		fmt.Printf("MagicStart0: 0x%08x\n", block.MagicStart0)
		fmt.Printf("MagicStart1: 0x%08x\n", block.MagicStart1)
		fmt.Printf("Flags: 0x%08x\n", block.Flags)
		fmt.Printf("TargetAddr: 0x%08x\n", block.TargetAddr)
		fmt.Printf("PayloadSize: 0x%08x\n", block.PayloadSize)
		fmt.Printf("BlockNo: 0x%08x\n", block.BlockNo)
		fmt.Printf("FileSize: 0x%08x\n", block.FileSize)
		fmt.Printf("\n\n")
		if err := binary.Write(writer, binary.LittleEndian, block); err != nil {
			return err
		}
	}

	return nil
}
