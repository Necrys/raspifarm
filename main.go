package main

import ( 
    "./bme280"
    "./config"
    "./trigger"
    "./sensorlog"
    "log"
    "fmt"
    "flag"
    "os"
    "os/signal"
    "syscall"
    "time"
    "github.com/stianeikeland/go-rpio"
)

func main() {
    isWorking := true
    sigs := make(chan os.Signal, 1)
    go func() {
        sig := <-sigs
        fmt.Printf("%v\n", sig)
        isWorking = false;
    }()

    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    // init GPIO lib
    err := rpio.Open()
    if err != nil {
        log.Fatal(err)
    }
    //defer rpio.Close()
    
    optCfgPath := flag.String("config", "default.json", "raspifarm configuration file")
    flag.Parse()

    cfg, err := config.Read(*optCfgPath)
    if err != nil {
        log.Fatal(err)
    }

    // TODO: do sensors creation by config
    conn, err := BME280.Connect(cfg.Sensors[0].Address, cfg.Sensors[0].Bus)
    if err != nil {
        log.Fatal(err)
    }

    defer conn.Disconnect()

    sensors := make(map[string]BME280.SensorIf)
    sensors[cfg.Sensors[0].Name] = conn

    ctx := trigger.CreateContext(cfg)
    
    err = sensorlog.InitFile(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // run sensor log
    go func() {
        logPeriod := time.Duration(cfg.Sensors_log_period) * time.Second // seconds
        for isWorking {
            err = sensorlog.WriteRow(ctx)
            if err != nil {
                log.Fatal(err)
            }
            time.Sleep(logPeriod)
        }
    }()

    // do hw test
    for k, v := range ctx.Relays {
        fmt.Printf("Testing relay \"%v\"\n", k)
        v.On()
        time.Sleep(time.Duration(500) * time.Millisecond)
        v.Off()
        time.Sleep(time.Duration(500) * time.Millisecond)
    }

    for isWorking {
        // Dirty hacks work only for ANSI terminals
        //fmt.Printf("\033[2J") // clear screen
        //fmt.Printf("\033[0;0H") // move cursor to up left corner

        // read sensors
        for _, s := range cfg.Sensors {
            temp, hum, pres, err := sensors[s.Name].ReadData()
            if err != nil {
                log.Fatal(err)
            }

            ctx.Sensors[s.Name].Temperature = temp
            ctx.Sensors[s.Name].Humidity = hum
            ctx.Sensors[s.Name].Pressure = pres

            //fmt.Println("----------------------------------------")
            //fmt.Printf("%v (%v)\n\r0n\r\r\rTemperature\t: %.2f C\nHumidity\t: %.2f %%\nPressure\t: %.2f mm Hg\n",
            //    s.Name, s.Type, ctx.Sensors[s.Name].Temperature, ctx.Sensors[s.Name].Humidity, ctx.Sensors[s.Name].Pressure)
            //fmt.Println("----------------------------------------")            
        }

        for _, trg := range ctx.Triggers {
            trigger.ProcessTrigger(trg, ctx)
        }

        time.Sleep(time.Duration(cfg.Update_period) * time.Millisecond)
    }

    // off all hw
    for _, v := range ctx.Relays {
        v.Off()
    }

    rpio.Close()
}
