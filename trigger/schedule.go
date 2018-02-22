package trigger

import (
    "../config"
    "errors"
    "log"
    "sort"
    "time"
)

type SchedulePoint struct {
    hour   int
    minute int
    action bool
    date   time.Time
}

type ScheduleTrigger struct {
    relay   string
    doAction bool
    points  []SchedulePoint
}

func NewScheduleTrigger(cfg config.Trigger) (*ScheduleTrigger, error) {
    this := &ScheduleTrigger {
        relay: cfg.Relay }

    now := time.Now()

    for _, p := range cfg.Time_points {
        newpoint := SchedulePoint{hour: p.Hour, minute: p.Minute}
        if p.Action == "On" {
            newpoint.action = true
        } else if p.Action == "Off" {
            newpoint.action = false
        } else {
            return nil, errors.New("Unknown action type for the trigger: \"" + p.Action + "\"")
        }

        newpoint.date = time.Date(now.Year(), now.Month(), now.Day(), newpoint.hour, newpoint.minute, 0, 0, now.Location())

        this.points = append(this.points, newpoint)
    }

    if len(this.points) == 0 {
        return nil, errors.New("No time points configured")
    }

    // sort points by time
    sort.Slice(this.points, func(i, j int) bool {
        if this.points[i].hour < this.points[j].hour {
            return true
        } else if (this.points[i].hour == this.points[j].hour) && (this.points[i].minute < this.points[j].minute) {
            return true
        }

        return false
    })

    // shift to current time period
    for {
        next := this.points[0].date

        if now.After(next) {
            // add a day and move to the end of list
            t := this.points[0]
            this.points = this.points[1:]

            t.date = t.date.Add(time.Hour * 24)
            this.points = append(this.points, t)
        } else {
            break
        }
    }

    return this, nil
}

func (this *ScheduleTrigger) Condition(ctx *Context) (bool) {
    if ctx.Relays[this.relay] == nil {
        log.Panic("Relay \"" + this.relay + "\" not found")
        return false
    }

    now := time.Now()

    next := this.points[0].date
    if now.After(next) {
        t := this.points[0]
        this.points = this.points[1:]

        t.date = t.date.Add(time.Hour * 24)
        this.points = append(this.points, t)
        this.doAction = t.action

        return true
    }

    return false
}

func (this *ScheduleTrigger) Action(ctx *Context) {
    if ctx.Relays[this.relay] == nil {
        log.Panic("Relay \"" + this.relay + "\" not found")
        return
    }

    if this.doAction == true {
        ctx.Relays[this.relay].On()
        ctx.log.Print("relay \"%s\" was switched on", this.relay)
    } else {
        ctx.Relays[this.relay].Off()
        ctx.log.Print("relay \"%s\" was switched off", this.relay)
    }
}