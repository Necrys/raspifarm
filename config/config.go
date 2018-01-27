package config

import (
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    Log_path string
    Sensors []Sensor
    Relays []Relay
    Triggers []Trigger
    Update_period uint32
    Sensors_log_period uint32
    Do_hw_test bool
}

type Sensor struct {
    Name string
    Type string
    Address uint8
    Bus int
}

type Relay struct {
    Name string
    Pin uint8
}

type Threshold struct {
    Value  float64
    Action string
}

type Trigger struct {
    Type           string
    Sensor         string
    Parameter      string
    Relay          string
    Low_threshold  Threshold
    High_threshold Threshold
}

func Read(cfgPath string) (*Config, error) {
    file, err := ioutil.ReadFile(cfgPath)
    if err != nil {
        return nil, err
    }

    var cfg Config
    err = json.Unmarshal(file, &cfg)
        if err != nil {
        return nil, err
    }

    return &cfg, nil
}
