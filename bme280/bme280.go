package BME280

import "fmt"
import "time"
import "golang.org/x/exp/io/i2c" // go get golang.org/x/exp/io/i2c

type Connection struct {
    conn *i2c.Device

    // settings
    oversampleHum byte
    oversampleTemp byte
    oversamplePres byte
    mode byte
}

// Registers (see datasheet at https://ae-bst.resource.bosch.com/media/_tech/media/datasheets/BST-BME280_DS001-12.pdf)
const (
    REG_CHIP_ID = 0xD0
    REG_CONTROL_HUM = 0xF2
    REG_CONTROL = 0xF4
    REG_CALIBRATION_01 = 0x88
    REG_CALIBRATION_02 = 0xA1
    REG_CALIBRATION_03 = 0xE1
    REG_DATA = 0xF7
)

// Initiate connection via I2C
func Connect(address uint8, bus int) (*Connection, error) {
    path := fmt.Sprintf("/dev/i2c-%d", bus)
    c, err := i2c.Open(&i2c.Devfs{Dev: path}, int(address))
    if err != nil {
        return nil, err
    }

    this := &Connection{conn: c}
    this.oversampleHum = 2
    this.oversampleTemp = 2
    this.oversamplePres = 2
    this.mode = 1

    return this, nil
}

// Close I2C connection
func (this *Connection) Disconnect() (error) {
    err := this.conn.Close()
    if err != nil {
        return err
    }
    return nil
}

// Read chip ID
func (this *Connection) ChipID() (byte, byte, error) {
    data := []byte{0, 0}
    err := this.conn.ReadReg(REG_CHIP_ID, data)
    if err != nil {
        return 0, 0, err
    }

    return data[0], data[1], nil
}

// Read temperature, humidity, pressure
func (this *Connection) ReadData() (float64, float64, float64, error) {
    // write control hum
    err := this.conn.WriteReg(REG_CONTROL_HUM, []byte{this.oversampleHum})
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

    // write other control
    err = this.conn.WriteReg(REG_CONTROL_HUM, []byte {this.oversampleTemp << 5 | this.oversamplePres << 2 | this.mode})
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

    // read calibration data
    calib01 := make([]byte, 24)
    calib02 := make([]byte, 1)
    calib03 := make([]byte, 7)

    err = this.conn.ReadReg(REG_CALIBRATION_01, calib01)
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

    err = this.conn.ReadReg(REG_CALIBRATION_02, calib02)
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

    err = this.conn.ReadReg(REG_CALIBRATION_03, calib03)
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

    digT1 := uint16(calib01[1]) << 8 | uint16(calib01[0])
    digT2 := int16(calib01[3]) << 8 | int16(calib01[2])
    digT3 := int16(calib01[5]) << 8 | int16(calib01[4])

    // wait for the measurements are done (Datasheet Appendix B: Measurement time and current calculation)
    waitTime := 1.25 + (2.3 * float64(this.oversampleTemp)) + ((2.3 * float64(this.oversamplePres)) + 0.575) + ((2.3 * float64(this.oversampleHum))+0.575);
    time.Sleep(time.Duration(waitTime) * time.Millisecond)

    // read measurements
    rawData := make([]byte, 8)
    err = this.conn.ReadReg(REG_DATA, rawData)
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

    rawTemp := uint32(rawData[3]) << 12 | uint32(rawData[4]) << 4 | uint32(rawData[5]) >> 4

    // refine temperature value
    var1 := ((uint32(rawTemp >> 3) - uint32(digT1) << 1) * uint32(digT2)) >> 11
    var2 := uint32(rawTemp >> 4) - uint32(digT1)
    var3 := (((var2 * var2) >> 12) * uint32(digT3)) >> 14
    vart := var1 + var3
    temperature := float64(((vart * 5) + 128) >> 8)

    return temperature/100.0, 0.0, 0.0, nil
}
