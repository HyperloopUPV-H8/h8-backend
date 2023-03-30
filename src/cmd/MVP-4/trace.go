package main

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func initTrace() *os.File {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	file, err := os.Create(*traceFile)
	if err != nil {
		trace.Logger = trace.Logger.Output(consoleWriter)
		trace.Fatal().Err(err).Msg("")
		return nil
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, file)

	global_logger := zerolog.New(multi).With().Timestamp().Caller().Logger()
	trace.Logger = global_logger

	if level, ok := traceLevelMap[*traceLevel]; ok {
		zerolog.SetGlobalLevel(level)
	} else {
		trace.Fatal().Msg("invalid log level selected")
		file.Close()
		return nil
	}

	return file
}
