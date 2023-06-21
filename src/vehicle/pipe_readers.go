package vehicle

import (
	"bufio"
	"fmt"
	"io"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
)

var pipeReaders map[uint16]common.ReaderFrom = map[uint16]common.ReaderFrom{
	1: NewDelimReaderFrom(0x00),
	2: NewDelimReaderFrom(0x00),
	3: NewDelimReaderFrom(0x00),
	4: NewDelimReaderFrom(0x00),
	5: NewStateOrderReaderFrom(),
	6: NewStateOrderReaderFrom(),
	7: NewStateSpaceReaderFrom(8, 15, 4),
}

func NewDelimReaderFrom(delim byte) DelimReaderFrom {
	return DelimReaderFrom{
		delim: delim,
	}
}

type DelimReaderFrom struct {
	delim byte
}

func (rf DelimReaderFrom) ReadFrom(r io.Reader) ([]byte, error) {
	reader := bufio.NewReader(r)
	buf, err := reader.ReadBytes(rf.delim)

	if err != nil {
		return buf, err
	}

	if len(buf) == 0 {
		return buf, nil
	}

	return buf[:len(buf)-1], nil
}

func NewStateOrderReaderFrom() StateOrderReaderFrom {
	return StateOrderReaderFrom{}
}

const BoardIdSizeLen = 1

type StateOrderReaderFrom struct{}

func (rf StateOrderReaderFrom) ReadFrom(r io.Reader) ([]byte, error) {
	reader := bufio.NewReader(r)
	idBoardIdAndSize, err := reader.Peek(BoardIdSizeLen)

	if err != nil {
		return nil, err
	}

	orderNum := idBoardIdAndSize[0]

	payload := make([]byte, BoardIdSizeLen+(orderNum*2))
	n, err := reader.Read(payload)

	if err != nil {
		return nil, err
	}

	if n != len(payload) {
		return nil, fmt.Errorf("expected %d bytes, got %d", len(payload), n)
	}

	return payload, nil
}

// 8*15 float32

type StateSpaceReaderFrom struct {
	rows         int
	cols         int
	variableSize int
}

func NewStateSpaceReaderFrom(rows int, cols int, varSize int) StateSpaceReaderFrom {
	return StateSpaceReaderFrom{
		rows:         rows,
		cols:         cols,
		variableSize: varSize,
	}
}

func (rf StateSpaceReaderFrom) ReadFrom(r io.Reader) ([]byte, error) {
	size := rf.cols * rf.rows * rf.variableSize
	stateSpaceBuf := make([]byte, size)
	n, err := r.Read(stateSpaceBuf)

	if err != nil {
		return nil, err
	}

	if n != size {
		return nil, fmt.Errorf("incorrect state space size: want %d got %d", size, n)
	}

	return stateSpaceBuf, nil
}
