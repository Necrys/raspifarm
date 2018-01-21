package BME280

import "fmt"
import "time"
import "log"
import "golang.org/x/exp/io/i2c" // go get golang.org/x/exp/io/i2c

type SensorIf interface {
    ReadData() (float64, float64, float64, error)
}

type Connection struct {
    conn *i2c.Device

    // settings
    oversampleHum byte
    oversampleTemp byte
    oversamplePres byte
    mode byte

    digT1 uint16
    digT2 int16
    digT3 int16

    digH1 uint8
    digH2 int16
    digH3 uint8
    digH4 int32
    digH5 int32
    digH6 int32

    digP1 uint16
    digP2 int16
    digP3 int16
    digP4 int16
    digP5 int16
    digP6 int16
    digP7 int16
    digP8 int16
    digP9 int16
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

    err = this.ReadCalibration()
    if err != nil {
        return nil, err
    }

    return this, nil
}

// Close I2C connection
func (this *Connection) Disconnect() () {
    err := this.conn.Close()
    if err != nil {
        log.Fatal(err)
    }
}

// Read calibration values
func (this *Connection) ReadCalibration() (error) {
    calib01 := make([]byte, 24)
    calib02 := make([]byte, 1)
    calib03 := make([]byte, 7)

    err := this.conn.ReadReg(REG_CALIBRATION_01, calib01)
    if err != nil {
        return err
    }

    err = this.conn.ReadReg(REG_CALIBRATION_02, calib02)
    if err != nil {
        return err
    }

    err = this.conn.ReadReg(REG_CALIBRATION_03, calib03)
    if err != nil {
        return err
    }

    this.digT1 = uint16(calib01[1]) << 8 | uint16(calib01[0])
    this.digT2 = int16(calib01[3]) << 8 | int16(calib01[2])
    this.digT3 = int16(calib01[5]) << 8 | int16(calib01[4])

    this.digH1 = uint8(calib02[0])
    this.digH2 = int16(calib03[1]) << 8 | int16(calib03[0])
    this.digH3 = uint8(calib03[2])

    this.digH4 = int32(calib03[3])
    this.digH4 = (this.digH4 << 24) >> 20
    this.digH4 = this.digH4 | (int32(calib03[4]) & 0x0F)

    this.digH5 = int32(calib03[5])
    this.digH5 = (this.digH5 << 24) >> 20
    this.digH5 = this.digH5 | (int32(calib03[4]) >> 4 & 0x0F)

    this.digH6 = int32(calib03[6])

    this.digP1 = uint16(calib01[7]) << 8 | uint16(calib01[6])
    this.digP2 = int16(calib01[9]) << 8 | int16(calib01[8])
    this.digP3 = int16(calib01[1]) << 8 | int16(calib01[10])
    this.digP4 = int16(calib01[13]) << 8 | int16(calib01[12])
    this.digP5 = int16(calib01[15]) << 8 | int16(calib01[14])
    this.digP6 = int16(calib01[17]) << 8 | int16(calib01[16])
    this.digP7 = int16(calib01[19]) << 8 | int16(calib01[18])
    this.digP8 = int16(calib01[21]) << 8 | int16(calib01[20])
    this.digP9 = int16(calib01[23]) << 8 | int16(calib01[22])

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
    control := []byte {this.oversampleTemp << 5 | this.oversamplePres << 2 | this.mode}
    err = this.conn.WriteReg(REG_CONTROL, control)
    if err != nil {
        return 0.0, 0.0, 0.0, err
    }

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
    rawHum  := uint32(rawData[6]) << 8  | uint32(rawData[7])
    rawPres := uint32(rawData[0]) << 12 | uint32(rawData[1]) << 4 | uint32(rawData[2]) >> 4

    // refine temperature value
    var1 := ((uint32(rawTemp >> 3) - uint32(this.digT1) << 1) * uint32(this.digT2)) >> 11
    var2 := uint32(rawTemp >> 4) - uint32(this.digT1)
    var3 := (((var2 * var2) >> 12) * uint32(this.digT3)) >> 14
    vart := var1 + var3
    temperature := float64(((vart * 5) + 128) >> 8)

    // refine pressure
    varp1 := float64(vart) / 2.0 - 64000.0
    varp2 := varp1 * varp1 * float64(this.digP6) / 32768.0
    varp2 = varp2 + varp1 * float64(this.digP5) * 2.0
    varp2 = varp2 / 4.0 + float64(this.digP4) * 65536.0
    varp1 = (float64(this.digP3) * varp1 * varp1 / 524288.0 + float64(this.digP2) * varp1) / 524288.0
    varp1 = (1.0 + varp1 / 32768.0) * float64(this.digP1)

    var pressure float64
    if varp1 == 0 {
        pressure=0
    } else {
        pressure = 1048576.0 - float64(rawPres)
        pressure = ((pressure - varp2 / 4096.0) * 6250.0) / varp1
        varp1 = float64(this.digP9) * pressure * pressure / 2147483648.0
        varp2 = pressure * float64(this.digP8) / 32768.0
        pressure = pressure + (varp1 + varp2 + float64(this.digP7)) / 16.0
        // convert to mmhg
        pressure = pressure / 1.33322387415
    }


    // refine humidity value
    humidity := float64(vart) - 76800.0
    humidity = (float64(rawHum) - (float64(this.digH4) * 64.0 + float64(this.digH5) / 16384.0 * humidity)) * (float64(this.digH2) / 65536.0 * (1.0 + float64(this.digH6) / 67108864.0 * humidity * (1.0 + float64(this.digH3) / 67108864.0 * humidity)))
    humidity = humidity * (1.0 - float64(this.digH1) * humidity / 524288.0)
    if humidity > 100.0 {
        humidity = 100.0
    } else if humidity < 0.0 {
        humidity = 0.0
    }

    return temperature/100.0, humidity, pressure/100.0, nil
}
