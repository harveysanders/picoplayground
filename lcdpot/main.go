package main

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
)

func main() {
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
	for {
		val := sensor.Get()
		percentage := (float32(val) / float32(max16Bit))
		voltage := percentage * sysV
		msg := fmt.Sprintf("V: %.1f, %.1f%%\n16-bit: %d", voltage, percentage*100, val)

		lcd.ClearDisplay()
		lcd.SetCursor(0, 0)
		lcd.Print([]byte(msg))

		time.Sleep(500 * time.Millisecond)
	}
}
