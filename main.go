package main

import ( 
    "./bme280"
    "./config"
    "log"
    "fmt"
    "flag"
)

func main() {
    optCfgPath := flag.String("config", "default.json", "raspifarm configuration file")
    flag.Parse()

    cfg, err := config.Read(*optCfgPath)
    if err != nil {
        log.Fatal(err)
    }

    conn, err := BME280.Connect(cfg.Sensors[0].Address, cfg.Sensors[0].Bus)
    if err != nil {
        log.Fatal(err)
    }

    id, ver, err := conn.ChipID()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Chip ID: %v\nChip version: %v\n", id, ver)

    err = conn.ReadCalibration()
    if err != nil {
        log.Fatal(err)
    }

    temp, hum, pres, err := conn.ReadData()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Temperature: %.2f C\nHumidity: %.2f %%\nPressure: %.2f mm Hg\n", temp, hum, pres)
    
    err = conn.Disconnect()
    if err != nil {
        log.Fatal(err)
    }
}
