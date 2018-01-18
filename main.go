package main

import ( 
    "./bme280"
    "log"
    "fmt"
)

func main() {
    conn, err := BME280.Connect(0x76, 1)
    if err != nil {
        log.Fatal(err)
    }

    id, ver, err := conn.ChipID()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Chip ID: %v\nChip version: %v\n", id, ver)
    
    err = conn.Disconnect()
    if err != nil {
        log.Fatal(err)
    }
}
