package actionlog

import (
    "os"
    "log"
    "bufio"
    "fmt"
)

type ActionLogIf interface {
    Print(string, ...interface{})
}

type actionLog struct {
    file *os.File
    writer *bufio.Writer
}

func NewActionLog(path string) (ActionLogIf, error) {
    f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Panic(err)
        return nil, err
    }

    this := &actionLog{file: f, writer: bufio.NewWriter(f)}
    return this, nil
}

func (this *actionLog) Print(format string, args ...interface{}) {
    fmt.Fprintf(this.writer, format + "\n", args)
}