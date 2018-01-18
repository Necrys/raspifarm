package main

import ( 
    "./bme280"
    "log"
)

func main() {
    conn, err := BME280.Connect(0x76, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    err = conn.Disconnect()
    if err != nil {
        log.Fatal(err)
    }
}
