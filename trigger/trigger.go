package trigger

import (
    "../config"
    "../relay"
    "log"
    "errors"
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
    Triggers []TriggerIf
}

type TriggerIf interface {
    Condition(ctx *Context) bool
    Action(ctx *Context)
}

func ProcessTrigger(trg TriggerIf, ctx *Context) {
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
            log.Panic(err)
            continue
        }
        this.Relays[r.Name] = rel
    }
    
    for _, t := range cfg.Triggers {
        switch t.Type {
        case "high_low_threshold": {
            trg, err := NewHighLowThresholdTrigger(t)
            if err != nil {
                log.Panic(err)
                continue
            }
            this.Triggers = append(this.Triggers, trg)
        }
        default:
            log.Panic("Unknown trigger type: \"" + t.Type + "\"")
            continue
        }
    }

    return this
}

func GetCtxSensorValue(ctx *Context, sensor string, param string) (float64, error) {
    if ctx.Sensors[sensor] == nil {
        return 0.0, errors.New("Sensor \"" + sensor + "\" not found")
    }
    
    switch param {
    case "temperature":
        return ctx.Sensors[sensor].Temperature, nil
    case "humidity":
        return ctx.Sensors[sensor].Humidity, nil
    case "pressure":
        return ctx.Sensors[sensor].Pressure, nil
    default:
        return 0.0, errors.New("Unknown parameter type: \"" + param + "\"")
    }
}