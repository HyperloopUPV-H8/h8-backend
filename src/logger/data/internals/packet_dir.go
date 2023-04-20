package internals

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/logger/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/data"
)

type PacketDir struct {
	files map[string]*models.SaveFile

	basePath string
}

func NewDir(basePath string) *PacketDir {
	return &PacketDir{
		files:    make(map[string]*models.SaveFile),
		basePath: basePath,
	}
}

func (dir *PacketDir) Write(meta packet.Metadata, packet data.Payload) (err error) {
	for name, value := range packet.Values {
		writeErr := dir.writeValues(name, meta, value)
		if writeErr != nil {
			err = writeErr
		}
	}

	return err
}

func (dir *PacketDir) writeValues(name string, meta packet.Metadata, value packet.Value) error {
	file, err := dir.getFile(name)
	if err != nil {
		return err
	}

	record := dir.toCSV(meta, value)

	return file.WriteCSV(record)
}

func (dir *PacketDir) toCSV(meta packet.Metadata, value packet.Value) []string {
	return []string{
		fmt.Sprint(meta.Timestamp.UnixNano()),
		fmt.Sprint(value.Inner()),
	}
}

func (dir *PacketDir) getFile(name string) (*models.SaveFile, error) {
	file, ok := dir.files[name]
	if !ok {
		var err error
		file, err = models.NewSaveFile(dir.basePath, name)
		if err != nil {
			return nil, err
		}
		dir.files[name] = file
	}

	return file, nil
}

func (dir *PacketDir) Flush() (err error) {
	for _, file := range dir.files {
		flushErr := file.Flush()
		if flushErr != nil {
			err = flushErr
		}
	}
	return err
}

func (dir *PacketDir) Close() (err error) {
	for _, file := range dir.files {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}
	return err
}
