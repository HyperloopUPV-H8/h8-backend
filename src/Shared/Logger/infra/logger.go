package infra

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/logger/domain"
)

type Logger struct {
	baseDir     string
	files       map[string]*os.File
	writeTicker *time.Ticker
	buffers     map[string]*[]string
	EntryChan   chan domain.Entry
	StopChan    chan bool
	Done        chan bool
}

func New(path string, delay time.Duration) *Logger {
	return newLogger(path, delay)
}

func newLogger(dir string, delay time.Duration) *Logger {
	baseDir := path.Join(dir, time.Now().Format("2006_01_02-15_04_05"))
	createLogDir(baseDir)
	return &Logger{
		baseDir:     baseDir,
		files:       make(map[string]*os.File),
		writeTicker: time.NewTicker(delay),
		buffers:     make(map[string]*[]string),
		EntryChan:   make(chan domain.Entry),
		StopChan:    make(chan bool),
		Done:        make(chan bool),
	}
}

func createLogDir(path string) {
	os.MkdirAll(path, 0777)
}

func (logger *Logger) AddEntry(entry domain.Entry) {
	logger.EntryChan <- entry
}

func (logger *Logger) addEntryToBuffer(entry domain.Entry) {
	buffer := logger.getBuffer(entry)
	bytesBuf, err := io.ReadAll(entry.Value)
	if err != nil {
		log.Fatalf("error reading entry reader: %v", err)
	}

	*buffer = append(*buffer, string(bytesBuf))
	logger.buffers[entry.Id] = buffer
}

func (logger *Logger) getBuffer(entry domain.Entry) *[]string {
	buffer, exists := logger.buffers[entry.Id]

	if !exists {
		logger.buffers[entry.Id] = new([]string)
		buffer = logger.buffers[entry.Id]
	}

	return buffer
}

func (logger *Logger) write() {
	for name := range logger.buffers {
		logger.writeBufferToFile(name)
	}
	logger.clearBuffers()
}

func (logger *Logger) writeBufferToFile(bufferName string) {
	file := logger.getFile(bufferName)
	for _, text := range *logger.buffers[bufferName] {
		file.WriteString(fmt.Sprintf("%v\n", text))
	}
}

func (logger *Logger) getFile(fileName string) *os.File {
	file, exists := logger.files[fileName]
	if !exists {
		file = createFile(logger.baseDir, fileName)
	}

	return file
}

func createFile(dirPath string, fileName string) *os.File {
	file, err := os.Create(path.Join(dirPath, fileName+".txt"))
	if err != nil {
		log.Fatalf("write values: %s\n", err)
	}

	return file
}

func (logger *Logger) clearBuffers() {
	logger.buffers = make(map[string]*[]string)
}

func (logger *Logger) Stop() {
	logger.StopChan <- true
	<-logger.Done
}

func (logger *Logger) close() {
	for _, file := range logger.files {
		file.Close()
	}
	logger.writeTicker.Stop()
	logger.Done <- true
}

func (logger *Logger) Record() {
	go func() {
	loop:
		for {
			select {
			case <-logger.writeTicker.C:
				logger.write()
			case entry := <-logger.EntryChan:
				logger.addEntryToBuffer(entry)
			case <-logger.StopChan:
				logger.write()
				logger.close()
				break loop
			}
		}
	}()
}
