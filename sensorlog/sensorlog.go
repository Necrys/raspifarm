package sensorlog

import (
    "../trigger"
    "fmt"
    "bufio"
    "time"
    "os"
    "errors"
)

var columnsCount int
var fileName string

func InitFile(ctx *trigger.Context) (error) {
    columnsCount = len(ctx.Sensors)

    t := time.Now()
    fileName = fmt.Sprintf("log_%d.%02d.%02dT%02d.%02d.%02d.csv", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second());
    file, err := os.Create(fileName)
    if err != nil {
        return err
    }
    
    defer file.Close()
    
    w := bufio.NewWriter(file)
    
    _, err = fmt.Fprintf(w, "Timestamp,")
    if err != nil {
        return err
    }

    c := 0
    for k, _ := range ctx.Sensors {
        _, err = fmt.Fprintf(w, "%s_temperature,%s_humidity,%s_pressure", k, k, k)
        if err != nil {
            return err
        }
        
        if c != columnsCount-1 {
            _, err = fmt.Fprintf(w, ",", k)
            if err != nil {
                return err
            }
        }
        
        c += 1
    }
    
    _, err = fmt.Fprintf(w, "\n")
    if err != nil {
        return err
    }
    
    w.Flush()

    return nil
}

func WriteRow(ctx *trigger.Context) (error) {
    if len(ctx.Sensors) != columnsCount {
        return errors.New("Incorrect sensors map size")
    }
    
    file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        return err
    }
    
    defer file.Close()
    
    w := bufio.NewWriter(file)

    t := time.Now()
    _, err = fmt.Fprintf(w, "%d.%02d.%02dT%02d.%02d.%02d,", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    if err != nil {
        return err
    }

    c := 0
    for _, v := range ctx.Sensors {
        _, err = fmt.Fprintf(w, "%v,%v,%v", v.Temperature, v.Humidity, v.Pressure)
        if err != nil {
            return err
        }

        if c != columnsCount-1 {
            _, err = fmt.Fprintf(w, ",")
            if err != nil {
                return err
            }
        }

        c += 1
    }
    
    _, err = fmt.Fprintf(w, "\n")
    if err != nil {
        return err
    }
    
    w.Flush()
    
    return nil
}
