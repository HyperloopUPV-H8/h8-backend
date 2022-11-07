package infra

import (
	"fmt"
	"log"
	"os"
	"path"
)

type FileData struct {
	handle *os.File
	values []ValueTimestamp
}

func NewFile(parentPath string, name string) *FileData {
	file, err := os.Create(path.Join(parentPath, fmt.Sprintf("%s.csv", name)))
	if err != nil {
		log.Fatalln("logger: new file:", err)
	}
	return &FileData{
		handle: file,
		values: make([]ValueTimestamp, 0),
	}
}

func (file *FileData) Dump() {
	for _, value := range file.values {
		file.write(value)
	}
	file.values = make([]ValueTimestamp, 0)
}

func (file *FileData) write(value ValueTimestamp) {
	file.handle.WriteString(fmt.Sprintf("%d, %s\n", value.timestamp.UnixNano(), value.value))
}

func (file *FileData) AddValue(value ValueTimestamp) {
	file.values = append(file.values, value)
}

type Dir struct {
	path  string
	files map[string]*FileData
}

func NewDir(base string) *Dir {
	os.MkdirAll(base, 0777)
	return &Dir{
		path:  base,
		files: make(map[string]*FileData),
	}
}

func (dir Dir) AddFile(valueName string) {
	dir.files[valueName] = NewFile(dir.path, valueName)
}

func (dir Dir) AppendValue(valueName string, value ValueTimestamp) bool {
	file, exists := dir.files[valueName]
	if !exists {
		dir.AddFile(valueName)
		file = dir.files[valueName]
	}

	file.AddValue(value)
	return true
}

func (dir Dir) Dump() {
	for _, file := range dir.files {
		file.Dump()
	}
}
