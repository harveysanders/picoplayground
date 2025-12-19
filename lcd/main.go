package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
)

func main() {
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

	lcd.ClearDisplay()
	lcd.SetCursor(0, 0)
	lcd.Print([]byte("Hello from TinyGo"))

	// Keep main() running
	for {
		println("done..")
		time.Sleep(time.Second * 5)
	}
}
