package trigger

import (
    "../config"
    "errors"
    "log"
)

type HighLowThresholdTrigger struct {
    sensor       string
    parameter    string
    relay        string
    low          float64
    high         float64
    last         float64
    actionOnLow  bool
    actionOnHigh bool
}

func NewHighLowThresholdTrigger(cfg config.Trigger) (*HighLowThresholdTrigger, error) {
    this := &HighLowThresholdTrigger {
        sensor: cfg.Sensor,
        parameter: cfg.Parameter,
        relay: cfg.Relay,
        low: cfg.Low_threshold.Value,
        high: cfg.High_threshold.Value,
        last: 0.0 }
    
    if cfg.Low_threshold.Action == "On" {
        this.actionOnLow = true
    } else if cfg.Low_threshold.Action == "Off" {
        this.actionOnLow = false
    } else {
        return nil, errors.New("Unknown action type for the trigger: \"" + cfg.Low_threshold.Action + "\"")
    }
    
    if cfg.High_threshold.Action == "On" {
        this.actionOnHigh = true
    } else if cfg.High_threshold.Action == "Off" {
        this.actionOnHigh = false
    } else {
        return nil, errors.New("Unknown action type for the trigger: \"" + cfg.High_threshold.Action + "\"")
    }
    
    return this, nil
}

func (this *HighLowThresholdTrigger) Condition(ctx *Context) (bool) {
    val, err := GetCtxSensorValue(ctx, this.sensor, this.parameter)
    if err != nil {
        log.Panic(err)
        return false
    }
        
    if this.last <= this.high && val > this.high {
        return true
    }
        
    if this.last >= this.low && val < this.low {
        return true
    }
    
    return false
}

func (this *HighLowThresholdTrigger) Action(ctx *Context) {
    val, err := GetCtxSensorValue(ctx, this.sensor, this.parameter)
    if err != nil {
        log.Panic(err)
        return
    }
    
    if ctx.Relays[this.relay] == nil {
        log.Panic("Relay \"" + this.relay + "\" not found")
        return
    }
    
    if this.last <= this.high && val > this.high {
        if this.actionOnHigh == true {
            ctx.Relays[this.relay].On()
        } else {
            ctx.Relays[this.relay].Off()
        }
    }
        
    if this.last >= this.low && val < this.low {
        if this.actionOnLow == true {
            ctx.Relays[this.relay].On()
        } else {
            ctx.Relays[this.relay].Off()
        }
    }

    this.last = val
}
