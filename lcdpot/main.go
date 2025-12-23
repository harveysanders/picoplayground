package main

import (
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
)

func main() {
	debugLED := machine.GP21
	debugLED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.InitADC()
	sensor := machine.ADC{Pin: machine.ADC0}
	sensor.Configure(machine.ADCConfig{})

	// Setup LCD display
	err := machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.GP4,
		SCL: machine.GP5,
	})
	if err != nil {
		for {
			println("could not configure I2C", err)
			time.Sleep(time.Second)
		}
	}

	// Try common addresses (0x27 then 0x3F)
	addrs := []uint8{0x27, 0x3F}
	var lcd hd44780i2c.Device
	found := false

	for _, a := range addrs {
		println("checking I2C address...")
		dev := hd44780i2c.New(machine.I2C0, a)
		dev.Configure(hd44780i2c.Config{
			Width:  16,
			Height: 2,
		})
		// Many drivers expose Connected()/Detect(); if not, just try writing.
		// (Check pkg.go.dev docs for the exact call set.)
		lcd = dev
		found = true
		break
	}

	if !found {
		for {
			println("LCD not found at 0x27/0x3F")
			time.Sleep(time.Second)
		}
	}

	var max16Bit uint16 = 65535
	var sysV float32 = 3.3

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

		debugLED.High()
		time.Sleep(250 * time.Millisecond)
		debugLED.Low()
		time.Sleep(250 * time.Millisecond)
	}
}
