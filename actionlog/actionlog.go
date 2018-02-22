package actionlog

import (
    "os"
    "log"
    "bufio"
    "fmt"
    "time"
)

type ActionLogIf interface {
    Print(string, ...interface{}) error
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

func (this *actionLog) Print(format string, args ...interface{}) error {
    t := time.Now()
    _, err := fmt.Fprintf(this.writer, "%d.%02d.%02dT%02d.%02d.%02d\t\t", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    if err != nil {
        return err
    }

    fmt.Fprintf(this.writer, format + "\n", args)
    if err != nil {
        return err
    }

    this.writer.Flush()

    return nil
}