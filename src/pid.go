package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	ErrProcessRunning = errors.New("process is running")
)

// Write writes a pidfile, returning an error
// if the process is already running or pidfile is orphaned
func WritePid(filename string) error {
	oldPid, err := getOldPid(filename)

	if err == nil && isProcessRunning(oldPid) {
		return ErrProcessRunning
	}

	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	buf := []byte(strconv.Itoa(os.Getpid()))
	_, err = file.Write(buf)

	if err != nil {
		return err
	}

	return nil
}

func getOldPid(filename string) (int, error) {
	oldPidFile, err := os.Open(filename)
	if err != nil {

		return 0, err
	}

	buf, err := io.ReadAll(oldPidFile)
	if err != nil {
		return 0, err
	}

	oldPid, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return 0, err
	}

	return oldPid, nil
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))

	if err != nil && err.Error() == "no such process" {
		return false
	}

	if err != nil && err.Error() == "os: process already finished" {
		return false
	}

	return true
}

// Remove a pidfile
func RemovePid(filename string) error {
	return os.RemoveAll(filename)
}
