package trigger

import (
    "../config"
    "../relay"
    "log"
)

type SensorData struct {
    Temperature float64
    Humidity float64
    Pressure float64
}

type RelayIf interface {
    State() bool
    On()
    Off()
}

type Context struct {
    Sensors map[string]*SensorData
    Relays map[string]RelayIf
}

type Trigger interface {
    Condition(ctx *Context) bool
    Action(ctx *Context)
}

func ProcessTrigger(trg Trigger, ctx *Context) {
    if trg.Condition(ctx) {
        trg.Action(ctx)
    }
}

func CreateContext(cfg *config.Config) (*Context) {
    this := &Context{Sensors: make (map[string]*SensorData), Relays: make (map[string]RelayIf)}

    for _, s := range cfg.Sensors {
        this.Sensors[s.Name] = &SensorData{}
    }
    
    for _, r := range cfg.Relays {
        rel, err := relay.NewRelay(r.Pin)
        if err != nil {
            log.Fatal(err)
            continue
        }
        this.Relays[r.Name] = rel
    }

    return this
}