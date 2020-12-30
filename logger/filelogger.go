// +build windows

package log

import (
	"log"
	"os"
	"path/filepath"

	"github.com/y-yagi/configure"
)

type Logger struct {
	file *os.File
	app  string
}

func NewLogger(app string) (*Logger, error) {
	dir := configure.ConfigDir(app)
	file := filepath.Join(dir, app+".log")
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	return &Logger{file: f, app: app}, nil
}

func (logger *Logger) Close() {
	logger.file.Close()
}

func (logger *Logger) Info(msg string) {
	log.Println("INFO:" + msg)
}

func (logger *Logger) Warning(msg string) {
	log.Println("WARNING:" + msg)
}

func (logger *Logger) Error(msg string) {
	log.Println("ERROR:" + msg)
}
