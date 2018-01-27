package relay

import (
    "github.com/stianeikeland/go-rpio"
)

type Relay struct {
    pin rpio.Pin
    state bool
}

func NewRelay(addr uint8) (*Relay, error) {
    this := &Relay{pin: rpio.Pin(addr), state: false}
    this.pin.Output()
    return this, nil
}

func (this* Relay) State() (bool) {
    return this.state
}

func (this *Relay) On() {
    if !this.state {
        this.pin.Low()
        this.state = !this.state
    }
}

func (this *Relay) Off() {
    if this.state {
        this.pin.High()
        this.state = !this.state
    }
}
