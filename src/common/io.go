package common

import "io"

func WriteAll(writer io.Writer, data []byte) (n int, err error) {
	written := 0
	for written < len(data) {
		n, err = writer.Write(data[written:])
		written += n

		if err != nil {
			break
		}
	}
	return written, err
}
