package BME280

import "fmt"
import "golang.org/x/exp/io/i2c" // go get golang.org/x/exp/io/i2c

type Connection struct {
    conn *i2c.Device
}

func Connect(address uint8, bus int) (*Connection, error) {
    path := fmt.Sprintf("/dev/i2c-%d", bus)
    c, err := i2c.Open(&i2c.Devfs{Dev: path}, int(address))
    if err != nil {
        return nil, err
    }

    this := &Connection{conn: c}
    return this, nil
}

func (this *Connection) Disconnect() (error) {
    err := this.conn.Close()
    if err != nil {
        return err
    }
    return nil
}

//
//    chipId := []byte{0}
//    err = conn.ReadReg(0xD0, chipId)
//    if err != nil {
//        log.Fatal(err)
//    }
//
//    fmt.Printf("%v\n", chipId[0])
//