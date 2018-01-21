package config

import (
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    Log_path string
    Sensors []Sensor
    Update_period uint32
}

type Sensor struct {
    Name string
    Type string
    Address uint8
    Bus int
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
