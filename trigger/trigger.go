package trigger

import (
    "../config"
)

type SensorData struct {
    Temperature float64
    Humidity float64
    Pressure float64
}

type SwitchIf interface {
    State() bool
    SwitchOn()
    SwitchOff()
}

type Context struct {
    Sensors map[string]*SensorData
    Switches map[string]SwitchIf
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
    this := &Context{Sensors: make (map[string]*SensorData), Switches: make (map[string]SwitchIf)}

    for _, s := range cfg.Sensors {
        this.Sensors[s.Name] = &SensorData{}
    }

    return this
}