package discord

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLog *log.Logger
)

type logWriter struct {
	f *os.File
}

func init() {
	infoLog = log.New(&logWriter{os.Stdout}, "", log.Lshortfile)
}

func (lw *logWriter) Write(b []byte) (int, error) {
	return fmt.Fprintf(lw.f, "%s: %s", time.Now().Format("3:04:05.000"), b)
}
