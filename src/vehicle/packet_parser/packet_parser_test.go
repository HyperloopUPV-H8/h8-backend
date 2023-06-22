package packet_parser

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestArrayParser(t *testing.T) {
	t.Run("encode array", func(t *testing.T) {
		prompt := []byte{0x00, 0x00, 0b0011, 0x00}
		want := []uint16{0, 3}

		r := bytes.NewReader(prompt)

		arr, err := readIntoArray[uint16](r, binary.LittleEndian, 2)

		if err != nil {
			t.Fatal(err)
		}

		typedArr, ok := arr.([]uint16)

		if !ok {
			t.Fatal("not correct type")
		}

		for index, item := range typedArr {
			if item != want[index] {
				t.Fatalf("want %d, got %d", want[index], item)
			}
		}
	})
}
