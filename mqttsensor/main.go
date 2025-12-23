package main

import (
	"errors"
	"log/slog"
	"machine"
	"strconv"
	"time"

	"github.com/harveysanders/picoplayground/mqttsensor/mqtt"
	"tinygo.org/x/drivers/hd44780i2c"
)

const (
	max16Bit      uint16  = 65535 // Max ADC value. The Pico has an onboard 16-bit ADC.
	sysV          float32 = 3.3   // Logic level in volts. Pico runs at 3.3VDC.
	serverAddrStr         = "10.0.0.9:1883"
)

func main() {
	start := time.Now()
	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	c := mqtt.Client{
		ID:         "tinygo-mqtt",
		Logger:     logger,
		Timeout:    5 * time.Second,
		TCPBufSize: 2030, // MTU - ethhdr - iphdr - tcphdr
	}

	// Buffered channel of 10 readings. We may need to adjust depending
	// on the sensor read frequency and network availability
	sensorReadings := make(chan mqtt.SensorReading, 10)
	go func() {
		err := c.ConnectAndPublish(serverAddrStr, sensorReadings)
		if err != nil {
			// Print error in a loop in case the serial monitor is not
			// ready before the inital messages
			printErrForever(logger, "connect to MQTT broker", slog.Any("reason", err))
		}
	}()

	debugLED := machine.GP21
	debugLED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.InitADC()
	sensor := machine.ADC{Pin: machine.ADC0}
	sensor.Configure(machine.ADCConfig{})

	// Setup LCD display over I2C
	err := machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.GP4,
		SCL: machine.GP5,
	})
	if err != nil {
		printErrForever(logger, "configure I2C", slog.Any("reason", err))
	}

	lcd, err := configureLCD(machine.I2C0)
	if err != nil {
		for {
			println(err.Error())
			time.Sleep(time.Second)
		}
	}

	lcd.ClearDisplay()
	// Buffer for LCD characters (16x2)
	// We need a preallocated buffer so the heap isn't exhausted
	// by many calls to fmt functions.
	printBuf := make([]byte, 0, 40)
	const floatNoExp = 'f'
	for {

		lcd.SetCursor(0, 0)

		val := sensor.Get()
		percentage := (float32(val) / float32(max16Bit))
		voltage := percentage * sysV
		// reslice the buffer to zero-length so append continues to work
		printBuf = printBuf[:0]
		printBuf = append(printBuf, "V: "...)
		printBuf = strconv.AppendFloat(printBuf, float64(voltage), floatNoExp, 1, 32)
		printBuf = append(printBuf, ", "...)
		printBuf = strconv.AppendFloat(printBuf, float64(percentage*100), floatNoExp, 1, 32)
		printBuf = append(printBuf, "%\n16-bit: "...)
		printBuf = strconv.AppendUint(printBuf, uint64(val), 10)

		lcd.Print(printBuf)

		sensorReadings <- mqtt.SensorReading{
			Voltage:   voltage,
			RawValue:  val,
			SinceBoot: time.Since(start),
		}

		debugLED.High()
		time.Sleep(250 * time.Millisecond)
		debugLED.Low()
		time.Sleep(250 * time.Millisecond)
	}
}

// configureLCD takes a preconfigured I2C peripheral and attempts to
// initialize the HD44780 LCD display. If no LCD found on the commond I2C
// addresses (0x27, 0x3F), an error is returned.
func configureLCD(i2c *machine.I2C) (hd44780i2c.Device, error) {
	// Try common addresses (0x27 then 0x3F)
	addrs := []uint8{0x27, 0x3F}
	var lcd hd44780i2c.Device
	found := false

	for _, a := range addrs {
		println("checking I2C address...")
		dev := hd44780i2c.New(i2c, a)
		dev.Configure(hd44780i2c.Config{
			Width:  16,
			Height: 2,
		})

		lcd = dev
		found = true
		break
	}

	if !found {
		return lcd, errors.New("LCD not found on addresses: 0x27, 0x3f")
	}
	return lcd, nil
}

// printError prints a string to serial @ 1hz. It
// blocks forever.
func printErrForever(logger *slog.Logger, msg string, args ...any) {
	for {
		logger.Error(msg, args...)
		time.Sleep(time.Second)
	}
}
