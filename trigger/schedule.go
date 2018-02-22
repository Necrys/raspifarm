package trigger

import (
    "../config"
    "errors"
    "log"
)

type SchedulePoint struct {
    hour   int
    minute int
    action bool
}

type ScheduleTrigger struct {
    relay   string
    initial bool
    points  []SchedulePoint
}

func NewScheduleTrigger(cfg config.Trigger) (*ScheduleTrigger, error) {
    this := &ScheduleTrigger {
        relay: cfg.relay
    }

    for _, p := range cfg.Time_points {
        
    }
    
    
    if cfg.hour < 0 || cfg.hour > 23 {
        return nil, errors.New("Incorrect time hour value")
    }

    if cfg.minute < 0 || cfg.minute > 59 {
        return nil, errors.New("Incorrect time minute value")
    }

    if cfg.Action == "On" {
        this.action = true
    } else if cfg.Action == "Off" {
        this.action = false
    } else {
        return nil, errors.New("Unknown action type for the trigger: \"" + cfg.Action + "\"")
    }
}